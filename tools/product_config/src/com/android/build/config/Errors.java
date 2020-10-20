/*
 * Copyright (C) 2020 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.android.build.config;

import java.lang.reflect.Field;
import java.io.PrintStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Reports errors to the caller.
 */
public class Errors {
    // Naming Convention:
    //   - ERROR_ for Categories with isLevelSettable false and Level.ERROR
    //   - WARNING_ for Categories with isLevelSettable false and default WARNING or HIDDEN
    //   - Don't have isLevelSettable true and not ERROR. (The constructor asserts this).
    public final Category ERROR_COMMAND_LINE = new Category(1, false, Level.ERROR,
            "Error on the command line.");
    public final Category WARNING_UNKNOWN_COMMAND_LINE_ERROR = new Category(2, true, Level.HIDDEN,
            "Passing unknown errors on the command line.  Hidden by default for\n"
            + "forward compatibilty.");
    public final Category REPLACE_ME = new Category(999, true, Level.WARNING,
            "When we have a default WARNING, replace this and update ErrorsTest to use that.");

    /**
     * The categories that are for this Errors object.
     */
    private final Map<Integer, Category> mCategories;

    /**
     * List of Entries that have occurred.
     */
    private final ArrayList<Entry> mEntries = new ArrayList();

    /**
     * Whether there has been a warning or an error yet.
     */
    private boolean mHadWarningOrError;

    /**
     * Whether there has been an error yet.
     */
    private boolean mHadError;

    /**
     * Whether errors are errors, warnings or hidden.
     */
    public static enum Level {
        HIDDEN("hidden"),
        WARNING("warning"),
        ERROR("error");

        private final String mLabel;

        Level(String label) {
            mLabel = label;
        }

        String getLabel() {
            return mLabel;
        }
    }

    /**
     * The available error codes.
     */
    public class Category {
        private final int mCode;
        private boolean mIsLevelSettable;
        private Level mLevel;
        private String mHelp;

        /**
         * Construct a Category object.
         */
        public Category(int code, boolean isLevelSettable, Level level, String help) {
            if (!isLevelSettable && level != Level.ERROR) {
                throw new RuntimeException("Don't have WARNING or HIDDEN without isLevelSettable");
            }
            mCode = code;
            mIsLevelSettable = isLevelSettable;
            mLevel = level;
            mHelp = help;
        }

        /**
         * Get the numeric code for the Category, which can be used to set the level.
         */
        public int getCode() {
            return mCode;
        }

        /**
         * Get whether the level of this Category can be changed.
         */
        public boolean isLevelSettable() {
            return mIsLevelSettable;
        }

        /**
         * Set the level of this category.
         */
        public void setLevel(Level level) {
            if (!mIsLevelSettable) {
                throw new RuntimeException("Can't set level for error " + mCode);
            }
            mLevel = level;
        }

        /**
         * Return the level, including any overrides.
         */
        public Level getLevel() {
            return mLevel;
        }

        /**
         * Return the category's help text.
         */
        public String getHelp() {
            return mHelp;
        }
    }

    /**
     * An instance of an error happening.
     */
    public class Entry {
        private final Category mCategory;
        private final Position mPosition;
        private final String mMessage;

        Entry(Category category, Position position, String message) {
            mCategory = category;
            mPosition = position;
            mMessage = message;
        }

        public Category getCategory() {
            return mCategory;
        }

        public Position getPosition() {
            return mPosition;
        }

        public String getMessage() {
            return mMessage;
        }
    }

    /**
     * Construct an Errors object. Note that for testing, there isn't a global Errors object.
     */
    public Errors() {
        HashMap<Integer, Category> categories = new HashMap();
        for (Field field: Errors.class.getFields()) {
            if (Category.class.equals(field.getType())) {
                Category category = null;
                try {
                    category = (Category)field.get(this);
                } catch (IllegalAccessException ex) {
                    // Wrap and rethrow, this is always on this class, so it's our programming
                    // error if this happens.
                    throw new RuntimeException("Categories on Errors should be public.", ex);
                }
                Category prev = categories.put(category.getCode(), category);
                if (prev != null) {
                    throw new RuntimeException("Duplicate categories with code "
                            + category.getCode());
                }
            }
        }
        mCategories = Collections.unmodifiableMap(categories);
    }

    /**
     * Returns a map of the category codes to the the categories.
     */
    public Map<Integer, Category> getCategories() {
        return mCategories;
    }

    /**
     * Add an error with no source position.
     */
    public void add(Category category, String message) {
        add(category, new Position(), message);
    }

    /**
     * Add an error.
     */
    public void add(Category category, Position pos, String message) {
        if (mCategories.get(category.getCode()) != category) {
            throw new RuntimeException("Errors.Category used from the wrong Errors object.");
        }
        mEntries.add(new Entry(category, pos, message));
        final Level level = category.getLevel();
        if (level == Level.WARNING || level == Level.ERROR) {
            mHadWarningOrError = true;
        }
        if (level == Level.ERROR) {
            mHadError = true;
        }
    }

    /**
     * Returns whether there has been a warning or an error yet.
     */
    public boolean hadWarningOrError() {
        return mHadWarningOrError;
    }

    /**
     * Returns whether there has been an error yet.
     */
    public boolean hadError() {
        return mHadError;
    }

    /**
     * Returns a list of all entries that were added.
     */
    public List<Entry> getEntries() {
        return new ArrayList<Entry>(mEntries);
    }

    /**
     * Prints the errors.
     */
    public void printErrors(PrintStream out) {
        for (Entry entry: mEntries) {
            final Category category = entry.getCategory();
            final Level level = category.getLevel();
            if (level == Level.HIDDEN) {
                continue;
            }
            out.println(entry.getPosition() + "[" + level.getLabel() + " "
                    + category.getCode() + "] " + entry.getMessage());
        }
    }
}
