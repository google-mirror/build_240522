/*
 * Copyright (C) 2008 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import java.util.regex.Pattern;
import java.util.regex.Matcher;
import java.util.ArrayList;

public class Comment
{
    static final Pattern LEADING_WHITESPACE = Pattern.compile(
                                "^[ \t\n\r]*(.*)$",
                                Pattern.DOTALL);

    static final Pattern TAG_BEGIN = Pattern.compile(
                                "[\r\n][\r\n \t]*@",
                                Pattern.DOTALL);

    static final Pattern TAG = Pattern.compile(
                                "(@[^ \t\r\n]+)[ \t\r\n]+(.*)",
                                Pattern.DOTALL);

    static final Pattern INLINE_TAG = Pattern.compile(
                                "(.*?)\\{(@[^ \t\r\n\\}]+)[ \t\r\n]*(.*?)\\}",
                                Pattern.DOTALL);

    static final Pattern FIRST_SENTENCE = Pattern.compile(
                                "((.*?)\\.)[ \t\r\n\\<](.*)",
                                Pattern.DOTALL);

    private static final String[] KNOWN_TAGS = new String[] {
            "@author",
            "@since",
            "@version",
            "@deprecated",
            "@undeprecate",
            "@docRoot",
            "@sdkCurrent",
            "@inheritDoc",
            "@more",
            "@code",
            "@samplecode",
            "@sample",
            "@include",
            "@serial",
        };

    public Comment(String text, ContainerInfo base, SourcePositionInfo sp)
    {
        mText = text;
        mBase = base;
        // sp now points to the end of the text, not the beginning!
        mPosition = SourcePositionInfo.findBeginning(sp, text);
    }

    private void parseRegex(String text)
    {
        Matcher m;

        m = LEADING_WHITESPACE.matcher(text);
        m.matches();
        text = m.group(1);

        m = TAG_BEGIN.matcher(text);

        int start = 0;
        int end = 0;
        while (m.find()) {
            end = m.start();

            tag(text, start, end);

            start = m.end()-1; // -1 is the @
        }
        end = text.length();
        tag(text, start, end);
    }

    private void tag(String text, int start, int end)
    {
        SourcePositionInfo pos = SourcePositionInfo.add(mPosition, mText, start);

        if (start >= 0 && end > 0 && (end-start) > 0) {
            text = text.substring(start, end);

            Matcher m = TAG.matcher(text);
            if (m.matches()) {
                // out of line tag
                tag(m.group(1), m.group(2), false, pos);
            } else {
                // look for inline tags
                m = INLINE_TAG.matcher(text);
                start = 0;
                while (m.find()) {
                    String str = m.group(1);
                    String tagname = m.group(2);
                    String tagvalue = m.group(3);
                    tag(null, m.group(1), true, pos);
                    tag(tagname, tagvalue, true, pos);
                    start = m.end();
                }
                int len = text.length();
                if (start != len) {
                    tag(null, text.substring(start), true, pos);
                }
            }
        }
    }

    private void tag(String name, String text, boolean isInline, SourcePositionInfo pos)
    {
        /*
        String s = isInline ? "inline" : "outofline";
        System.out.println("---> " + s
                + " name=[" + name + "] text=[" + text + "]");
        */
        if (name == null) {
            mInlineTagsList.add(new TextTagInfo("Text", "Text", text, pos));
        }
        else if (name.equals("@param")) {
            mParamTagsList.add(new ParamTagInfo("@param", "@param", text, mBase, pos));
        }
        else if (name.equals("@see")) {
            mSeeTagsList.add(new SeeTagInfo("@see", "@see", text, mBase, pos));
        }
        else if (name.equals("@link") || name.equals("@linkplain")) {
            mInlineTagsList.add(new SeeTagInfo(name, "@see", text, mBase, pos));
        }
        else if (name.equals("@throws") || name.equals("@exception")) {
            mThrowsTagsList.add(new ThrowsTagInfo("@throws", "@throws", text, mBase, pos));
        }
        else if (name.equals("@return")) {
            mReturnTagsList.add(new ParsedTagInfo("@return", "@return", text, mBase, pos));
        }
        else if (name.equals("@deprecated")) {
            if (text.length() == 0) {
                Errors.error(Errors.MISSING_COMMENT, pos,
                        "@deprecated tag with no explanatory comment");
                text = "No replacement.";
            }
            mDeprecatedTagsList.add(new ParsedTagInfo("@deprecated", "@deprecated", text, mBase, pos));
        }
        else if (name.equals("@literal")) {
            mInlineTagsList.add(new LiteralTagInfo(name, name, text, pos));
        }
        else if (name.equals("@hide") || name.equals("@pending") || name.equals("@doconly")) {
            // nothing
        }
        else if (name.equals("@attr")) {
            AttrTagInfo tag = new AttrTagInfo("@attr", "@attr", text, mBase, pos);
            mAttrTagsList.add(tag);
            Comment c = tag.description();
            if (c != null) {
                for (TagInfo t: c.tags()) {
                    mInlineTagsList.add(t);
                }
            }
        }
        else if (name.equals("@undeprecate")) {
            mUndeprecateTagsList.add(new TextTagInfo("@undeprecate", "@undeprecate", text, pos));
        }
        else if (name.equals("@include") || name.equals("@sample")) {
            mInlineTagsList.add(new SampleTagInfo(name, "@include", text, mBase, pos));
        }
        else {
            boolean known = false;
            for (String s: KNOWN_TAGS) {
                if (s.equals(name)) {
                    known = true;
                    break;
                }
            }
            if (!known) {
                known = DroidDoc.knownTags.contains(name);
            }
            if (!known) {
                Errors.error(Errors.UNKNOWN_TAG, pos == null ? null : new SourcePositionInfo(pos),
                        "Unknown tag: " + name);
            }
            TagInfo t = new TextTagInfo(name, name, text, pos);
            if (isInline) {
                mInlineTagsList.add(t);
            } else {
                mTagsList.add(t);
            }
        }
    }

    private void parseBriefTags()
    {
        int N = mInlineTagsList.size();

        // look for "@more" tag, which means that we might go past the first sentence.
        int more = -1;
        for (int i=0; i<N; i++) {
            if (mInlineTagsList.get(i).name().equals("@more")) {
                more = i;
            }
        }
          if (more >= 0) {
            for (int i=0; i<more; i++) {
                mBriefTagsList.add(mInlineTagsList.get(i));
            }
        } else {
            for (int i=0; i<N; i++) {
                TagInfo t = mInlineTagsList.get(i);
                if (t.name().equals("Text")) {
                    Matcher m = FIRST_SENTENCE.matcher(t.text());
                    if (m.matches()) {
                        String text = m.group(1);
                        TagInfo firstSentenceTag = new TagInfo(t.name(), t.kind(), text, t.position());
                        mBriefTagsList.add(firstSentenceTag);
                        break;
                    }
                }
                mBriefTagsList.add(t);

            }
        }
    }

    public TagInfo[] tags()
    {
        init();
        return mInlineTags;
    }

    public TagInfo[] tags(String name)
    {
        init();
        ArrayList<TagInfo> results = new ArrayList<TagInfo>();
        int N = mInlineTagsList.size();
        for (int i=0; i<N; i++) {
            TagInfo t = mInlineTagsList.get(i);
            if (t.name().equals(name)) {
                results.add(t);
            }
        }
        return results.toArray(new TagInfo[results.size()]);
    }

    public ParamTagInfo[] paramTags()
    {
        init();
        return mParamTags;
    }

    public SeeTagInfo[] seeTags()
    {
        init();
        return mSeeTags;
    }

    public ThrowsTagInfo[] throwsTags()
    {
        init();
        return mThrowsTags;
    }

    public TagInfo[] returnTags()
    {
        init();
        return mReturnTags;
    }

    public TagInfo[] deprecatedTags()
    {
        init();
        return mDeprecatedTags;
    }

    public TagInfo[] undeprecateTags()
    {
        init();
        return mUndeprecateTags;
    }

    public AttrTagInfo[] attrTags()
    {
        init();
        return mAttrTags;
    }

    public TagInfo[] briefTags()
    {
        init();
        return mBriefTags;
    }

    // This method is used to replace indexOf usage below which seems to cause
    // some problems with the Java environment.  See the comment in isHidden
    // to understand the background.
    private static boolean isFoundIn(String needle, String haystack)
    {
	char[] n = needle.toCharArray();
	char[] h = haystack.toCharArray();
	int last = haystack.length() - needle.length();
	int nlen = needle.length();
	int i;
	int j;
	boolean match = false;

	for (i = 0; !match && i < last; i++) {
	    j = 0;
	    if (h[i] == n[j]) {
		match = true;
		for (; j < nlen; j++) {
		    if (h[i+j] != n[j]) {
			match = false;
			break;
		    }
		}
	    }
	}
	return match;
    }


    public boolean isHidden()
    {
        if (mHidden >= 0) {
            return mHidden != 0;
        } else {
            if (DroidDoc.checkLevel(DroidDoc.SHOW_HIDDEN)) {
                mHidden = 0;
                return false;
            }
	    // The original code from droiddoc had this perfectly reasonable use
	    // of indexOf.
            // boolean b = mText.indexOf("@hide") >= 0 || mText.indexOf("@pending") >= 0;
            // mHidden = b ? 1 : 0;
	    // if (b) {
	    // Unfortunately under some conditions, this code mysteriously returns
	    // true for an mText that does not contain @hide or @pending.
	    // This seems to apply with the Intent class and cause
	    // document generation to fail.  This failure was not reliable and could
	    // only be reproduced on (between) 50-150 builds.  The change to use
	    // the new method isFoundIn is stable over a much higher number of builds
	    // although the true root cause is not yet discovered.
	    if (isFoundIn("@hide", mText) || isFoundIn("@pending", mText)) {
		mHidden = 1;
		return true;
	    } else {
		mHidden = 0;
		return false;
	    }
            //return b;
        }
    }

    public boolean isDocOnly() {
        if (mDocOnly >= 0) {
            return mDocOnly != 0;
        } else {
            boolean b = (mText != null) && (mText.indexOf("@doconly") >= 0);
            mDocOnly = b ? 1 : 0;
            return b;
        }
    }

    private void init()
    {
        if (!mInitialized) {
            initImpl();
        }
    }

    private void initImpl()
    {
        isHidden();
        isDocOnly();

        // Don't bother parsing text if we aren't generating documentation.
        if (DroidDoc.parseComments()) {
            parseRegex(mText);
            parseBriefTags();
        } else {
          // Forces methods to be recognized by findOverriddenMethods in MethodInfo.
          mInlineTagsList.add(new TextTagInfo("Text", "Text", mText,
                  SourcePositionInfo.add(mPosition, mText, 0)));
        }

        mText = null;
        mInitialized = true;

        mInlineTags = mInlineTagsList.toArray(new TagInfo[mInlineTagsList.size()]);
        mParamTags = mParamTagsList.toArray(new ParamTagInfo[mParamTagsList.size()]);
        mSeeTags = mSeeTagsList.toArray(new SeeTagInfo[mSeeTagsList.size()]);
        mThrowsTags = mThrowsTagsList.toArray(new ThrowsTagInfo[mThrowsTagsList.size()]);
        mReturnTags = ParsedTagInfo.joinTags(mReturnTagsList.toArray(
                                             new ParsedTagInfo[mReturnTagsList.size()]));
        mDeprecatedTags = ParsedTagInfo.joinTags(mDeprecatedTagsList.toArray(
                                        new ParsedTagInfo[mDeprecatedTagsList.size()]));
        mUndeprecateTags = mUndeprecateTagsList.toArray(new TagInfo[mUndeprecateTagsList.size()]);
        mAttrTags = mAttrTagsList.toArray(new AttrTagInfo[mAttrTagsList.size()]);
        mBriefTags = mBriefTagsList.toArray(new TagInfo[mBriefTagsList.size()]);

        mParamTagsList = null;
        mSeeTagsList = null;
        mThrowsTagsList = null;
        mReturnTagsList = null;
        mDeprecatedTagsList = null;
        mUndeprecateTagsList = null;
        mAttrTagsList = null;
        mBriefTagsList = null;
    }

    boolean mInitialized;
    int mHidden = -1;
    int mDocOnly = -1;
    String mText;
    ContainerInfo mBase;
    SourcePositionInfo mPosition;
    int mLine = 1;

    TagInfo[] mInlineTags;
    TagInfo[] mTags;
    ParamTagInfo[] mParamTags;
    SeeTagInfo[] mSeeTags;
    ThrowsTagInfo[] mThrowsTags;
    TagInfo[] mBriefTags;
    TagInfo[] mReturnTags;
    TagInfo[] mDeprecatedTags;
    TagInfo[] mUndeprecateTags;
    AttrTagInfo[] mAttrTags;

    ArrayList<TagInfo> mInlineTagsList = new ArrayList<TagInfo>();
    ArrayList<TagInfo> mTagsList = new ArrayList<TagInfo>();
    ArrayList<ParamTagInfo> mParamTagsList = new ArrayList<ParamTagInfo>();
    ArrayList<SeeTagInfo> mSeeTagsList = new ArrayList<SeeTagInfo>();
    ArrayList<ThrowsTagInfo> mThrowsTagsList = new ArrayList<ThrowsTagInfo>();
    ArrayList<TagInfo> mBriefTagsList = new ArrayList<TagInfo>();
    ArrayList<ParsedTagInfo> mReturnTagsList = new ArrayList<ParsedTagInfo>();
    ArrayList<ParsedTagInfo> mDeprecatedTagsList = new ArrayList<ParsedTagInfo>();
    ArrayList<TagInfo> mUndeprecateTagsList = new ArrayList<TagInfo>();
    ArrayList<AttrTagInfo> mAttrTagsList = new ArrayList<AttrTagInfo>();


}
