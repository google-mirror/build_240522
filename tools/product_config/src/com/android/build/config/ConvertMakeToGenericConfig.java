/*
 * Copyright (C) 2021 The Android Open Source Project
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

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * Converts a MakeConfig into a Generic config by applying heuristics about
 * the types of variable assignments that we do.
 */
public class ConvertMakeToGenericConfig {
    public static GenericConfig convert(MakeConfig make) {
        final GenericConfig result = new GenericConfig();

        // Base class fields
        result.setPhase(make.getPhase());
        result.setRootNodes(make.getRootNodes());
        for (Map.Entry<String, ConfigBase.VarType> entry: make.getProductVars().entrySet()) {
            result.addProductVar(entry.getKey(), entry.getValue());
        }

        // Each file
        for (MakeConfig.ConfigFile f: make.getConfigFiles()) {
            final GenericConfig.ConfigFile genericFile
                    = new GenericConfig.ConfigFile(f.getFilename());
            result.addConfigFile(genericFile);

            final List<MakeConfig.Block> blocks = f.getBlocks();

            // Some assertions:
            // TODO: Include better context for these errors.
            // There should always be at least a BEGIN and an AFTER, so assert this.
            if (blocks.size() < 2) {
                throw new RuntimeException("expected at least blocks.size() >= 2. Actcual size: "
                        + blocks.size());
            }
            if (blocks.get(0).getBlockType() != MakeConfig.BlockType.BEFORE) {
                throw new RuntimeException("expected first block to be BEFORE");
            }
            if (blocks.get(blocks.size() - 1).getBlockType() != MakeConfig.BlockType.AFTER) {
                throw new RuntimeException("expected first block to be AFTER");
            }
            // Everything in between should be an INHERIT block.
            for (int index = 1; index < blocks.size() - 1; index++) {
                if (blocks.get(index).getBlockType() != MakeConfig.BlockType.INHERIT) {
                    throw new RuntimeException("expected INHERIT at block " + index);
                }
            }

            // Each block represents a snapshot of the interpreter variable state (minus a few big
            // sets of variables which we don't export because they're used in the internals
            // of node_fns.mk, so we know they're not necessary here). The first (BEFORE) one
            // is everything that is set before the file is included, so it forms the base
            // for everything else.
            MakeConfig.Block prevBlock = blocks.get(0);

            for (int index = 1; index < blocks.size(); index++) {
                final MakeConfig.Block block = blocks.get(index);
                for (final Map.Entry<String, Str> entry: block.getVars().entrySet()) {
                    final String varName = entry.getKey();
                    final Str varVal = entry.getValue();
                    final Str prevVal = prevBlock.getVar(varName);

                    if (prevVal == null) {
                        // New variable.
                        genericFile.addStatement(new GenericConfig.Assign(varName, varVal));
                    } else if (!varVal.equals(prevVal)) {
                        // The value changed from the last block.
                        if (varVal.equals("")) {
                            // It was set to empty
                            genericFile.addStatement(new GenericConfig.Assign(varName, varVal));
                        } else {
                            // Product vars have the @inherit processing. Other vars we
                            // will just ignore and put in one section at the end, based
                            // on the difference between the BEFORE and AFTER blocks.
                            switch (make.getVarType(varName)) {
                                case LIST:
                                case SINGLE:
                                    handleListVariable(genericFile, block, varName, varVal,
                                            prevVal);
                                    break;
                                case UNKNOWN:
                                    if (block.getBlockType() == MakeConfig.BlockType.AFTER) {
                                        // For UNKNOWN variables, we don't worry about the
                                        // intermediate steps, just take the final value.
                                        genericFile.addStatement(
                                                new GenericConfig.Assign(varName, varVal));
                                    }
                                    break;
                            }
                        }
                    }
                }
                // Handle variables that are in prevBlock but not block -- they were
                // deleted. Is this even possible, or do they show up as ""?  We will
                // treat them as positive assigments to empty string
                for (String prevName: prevBlock.getVars().keySet()) {
                    if (!block.getVars().containsKey(prevName)) {
                        genericFile.addStatement(
                                new GenericConfig.Assign(prevName, new Str("")));
                    }
                }
                if (block.getBlockType() == MakeConfig.BlockType.INHERIT) {
                    genericFile.addStatement(
                            new GenericConfig.Inherit(block.getInheritedFile()));
                }
                // For next iteration
                prevBlock = block;
            }
        }
        return result;
    }

    /**
     * Handle the special inherited list values, where the inherit-product puts in the
     * @inherit:... markers, adding Statements to the ConfigFile.
     */
    static void handleListVariable(GenericConfig.ConfigFile genericFile,
            MakeConfig.Block block, String varName, Str varVal, Str prevVal) {
        String varText = varVal.toString();
        String prevText = prevVal.toString().trim();
        if (block.getBlockType() == MakeConfig.BlockType.INHERIT) {
            // inherit-product appends @inherit:... so drop that.
            final String marker =
                    "@inherit:" + block.getInheritedFile();
            if (varText.endsWith(marker)) {
                varText = varText.substring(0, varText.length() - marker.length()).trim();
            } else {
                mErrors.ERROR_IMPROPER_PRODUCT_VAR_MARKER.add(varVal.getPosition,
                        "Variable didn't end with marker \"" + marker + "\": " + varText);
            }
        }

        if (!varText.equals(prevText)) {
            // If the variable value was actually changed.
            final ArrayList<String> words = split(varText, prevText);
            if (words.size() == 0) {
                // Pure Assignment, none of the previous value is present.
                genericFile.addStatement(
                        new GenericConfig.Assign(varName, new Str(varVal.getPosition(), varText)));
            } else {
                // Self referential value (prepend, append, both).
                genericFile.addStatement(
                        new GenericConfig.Assign(varName, Str.toList(varVal.getPosition(), words)));
                if (words.size() > 2) {
                    // This is indicative of a construction that might not be quite
                    // what we want.  The above code will do something that works if it was
                    // of the form "VAR := a $(VAR) b $(VAR) c", but if the original code
                    // something else this won't work. This doesn't happen in AOSP, but
                    // it's a theoretically possiblity, so someone might do it.
                    mErrors.WARNING_VARIABLE_RECURSION.add(
                            new Position(words.get(2).getPosition()),
                            "Possible unsupported variable recursion: "
                                + varName + " = " + varVal + " (prev=" + prevVal + ")");
                }
            }
        }
        // else Variable not touched
    }

    /**
     * Split 'haystack' on occurrences of 'needle'. Trims each string of whitespace
     * to preserve make list semantics.
     */
    static ArrayList<String> split(String haystack, String needle) {
        final ArrayList<String> result = new ArrayList();
        final int needleLen = needle.length();
        if (needleLen == 0) {
            return result;
        }
        int start = 0;
        int end;
        while ((end = haystack.indexOf(needle, start)) >= 0) {
            result.add(haystack.substring(start, end).trim());
            start = end + needleLen;
        }
        result.add(haystack.substring(start).trim());
        return result;
    }
}
