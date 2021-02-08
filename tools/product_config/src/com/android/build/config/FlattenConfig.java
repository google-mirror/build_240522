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

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import java.util.regex.Pattern;

public class FlattenConfig {
    private static final Pattern RE_SPACE = Pattern.compile("\\p{Space}+");

    private final Errors mErrors;
    private final GenericConfig mGenericConfig;
    private final Map<String, GenericConfig.ConfigFile> mGenericConfigs;
    private final FlatConfig mResult = new FlatConfig();
    private final TreeMap<String, FlatConfig.Value> mVariables;
    /**
     * Files that have been visited, to prevent infinite recursion. There are no
     * conditionals at this point in the processing, so we don't need a stack, just
     * a single set.
     */
    private final Set<String> mVisitedFiles = new HashSet();


    private FlattenConfig(Errors errors, GenericConfig genericConfig) {
        mErrors = errors;
        mGenericConfig = genericConfig;
        mGenericConfigs = genericConfig.getFiles();
        mVariables = mResult.getValues();

        // Base class fields
        mResult.copyFrom(genericConfig);
    }

    /**
     * Flatten a GenericConfig to a FlatConfig, with 'root' as the top-level file
     * where PRODUCT_NAME is TARGET_PRODUCT.
     *
     * Makes three passes through the genericConfig, one to flatten the single variables,
     * one to flatten the list variables, and one to flatten the unknown variables. Each
     * has a slightly different algorithm.
     */
    public static FlatConfig flatten(Errors errors, GenericConfig genericConfig, Str root) {
        final FlattenConfig flattener = new FlattenConfig(errors, genericConfig);

        // TODO: Do we need to worry about the initial state of variables? Anything
        // that from the product config

        flattener.flattenListVars(root);
        flattener.flattenSingleVars(root);
        flattener.flattenUnknownVars(root);

        return flattener.mResult;
    }

    interface AssignCallback {
        void onAssignStatement(GenericConfig.Assign assign);
    }

    interface InheritCallback {
        void onInheritStatement(GenericConfig.Inherit assign);
    }

    /**
     * Do a bunch of validity checks, and then iterate through each of the statements
     * in the given file.  For Assignments, the callback is only called for variables
     * matching varType.
     */
    private void forEachStatement(Str filename, ConfigBase.VarType varType,
            AssignCallback assigner, InheritCallback inheriter) {
        if (mVisitedFiles.contains(filename.toString())) {
            mErrors.ERROR_INFINITE_RECURSION.add(filename.getPosition(),
                    "File is already in the inherit-product stack: " + filename);
            return;
        }

        mVisitedFiles.add(filename.toString());
        try {
            final GenericConfig.ConfigFile genericFile = mGenericConfigs.get(filename.toString());

            if (genericFile == null) {
                mErrors.ERROR_MISSING_CONFIG_FILE.add(filename.getPosition(),
                        "Unable to find config file: " + filename);
                return;
            }

            for (final GenericConfig.Statement statement: genericFile.getStatements()) {
                if (statement instanceof GenericConfig.Assign) {
                    if (assigner != null) {
                        final GenericConfig.Assign assign = (GenericConfig.Assign)statement;
                        final String varName = assign.getName();

                        // Assert that we're not stomping on another variable, which
                        // really should be impossible at this point.
                        assertVarType(filename.toString(), varName);

                        if (mGenericConfig.getVarType(varName) == varType) {
                            assigner.onAssignStatement(assign);
                        }
                    }
                } else if (statement instanceof GenericConfig.Inherit) {
                    if (inheriter != null) {
                        inheriter.onInheritStatement((GenericConfig.Inherit)statement);
                    }
                }
            }
        } finally {
            // Also executes after return statements, so we always remove this.
            mVisitedFiles.remove(filename);
        }
    }

    /**
     * Traverse the inheritance hierarchy, setting list-value product config variables.
     */
    private void flattenListVars(final Str filename) {
        forEachStatement(filename, ConfigBase.VarType.LIST,
                (assign) -> {
                    // Append to the list
                    appendToListVar(assign.getName(), assign.getValue());
                    // TODO: Drop 2nd ones for diamond inheritance
                },
                (inherit) -> {
                    flattenListVars(inherit.getFilename());
                });
    }

    /**
     * Traverse the inheritance hierarchy, setting single-value product config variables.
     */
    private void flattenSingleVars(final Str filename) {
        // flattenSingleVars has two loops.  The first sets all variables that are
        // defined for *this* file.  The second traverses through the inheritance,
        // to fill in values that weren't defined in this file.  The first appearance of
        // the variable is the one that wins.

        forEachStatement(filename, ConfigBase.VarType.SINGLE,
                (assign) -> {
                    final String varName = assign.getName();
                    FlatConfig.Value v = mVariables.get(varName);
                    // Only take the first value that we see for single variables.
                    FlatConfig.Value value = mVariables.get(varName);
                    if (!mVariables.containsKey(varName)) {
                        final List<Str> valueList = assign.getValue();
                        // There should never be more than one item in this list, because
                        // SINGLE values should never be appended to.
                        if (valueList.size() != 1) {
                            final StringBuilder positions = new StringBuilder("[");
                            for (Str s: valueList) {
                                positions.append(s.getPosition());
                            }
                            positions.append(" ]");
                            throw new RuntimeException("Value list found for SINGLE variable "
                                    + varName + " size=" + valueList.size()
                                    + "positions=" + positions.toString());
                        }
                        mVariables.put(varName,
                                new FlatConfig.Value(ConfigBase.VarType.SINGLE,
                                    valueList.get(0)));
                    }
                }, null);

        forEachStatement(filename, ConfigBase.VarType.SINGLE, null,
                (inherit) -> {
                    flattenSingleVars(inherit.getFilename());
                });
    }

    /**
     * Traverse the inheritance hierarchy and flatten the values 
     */
    private void flattenUnknownVars(Str filename) {
        // flattenUnknownVars has two loops: First to attempt to set the variable from
        // this file, and then a second loop to handle the inheritance.  This is odd
        // but it matches the order the files are included in node_fns.mk. The last appearance
        // of the value is the one that wins.

        forEachStatement(filename, ConfigBase.VarType.UNKNOWN,
                (assign) -> {
                    // Overwrite the current value with whatever is now in the file.
                    mVariables.put(assign.getName(),
                            new FlatConfig.Value(ConfigBase.VarType.UNKNOWN,
                                flattenAssignList(assign, new Str(""))));
                }, null);

        forEachStatement(filename, ConfigBase.VarType.UNKNOWN, null,
                (inherit) -> {
                    flattenUnknownVars(inherit.getFilename());
                });
    }

    /**
     * Throw an exception if there's an existing variable with a different type.
     */
    private void assertVarType(String filename, String varName) {
        if (mGenericConfig.getVarType(varName) == ConfigBase.VarType.UNKNOWN) {
            final FlatConfig.Value prevValue = mVariables.get(varName);
            if (prevValue != null
                    && prevValue.getVarType() != ConfigBase.VarType.UNKNOWN) {
                throw new RuntimeException("Mismatched var types:"
                        + " filename=" + filename
                        + " varType=" + mGenericConfig.getVarType(varName)
                        + " varName=" + varName
                        + " prevValue=" + prevValue.debugString());
            }
        }
    }

    /**
     * Appends all of the words in in 'items' to an entry in mVariables keyed by 'varName',
     * creating one if necessary.
     */
    private void appendToListVar(String varName, List<Str> items) {
        FlatConfig.Value value = mVariables.get(varName);
        if (value == null) {
            value = new FlatConfig.Value(new ArrayList());
            mVariables.put(varName, value);
        }
        final List<Str> out = value.getList();
        for (Str item: items) {
            for (String word: RE_SPACE.split(item.toString())) {
                if (word.length() > 0) {
                    out.add(new Str(item.getPosition(), word));
                }
            }
        }
    }

    /**
     * Flatten the list of strings in an Assign statement, using the previous value
     * as a separator.
     */
    private Str flattenAssignList(GenericConfig.Assign assign, Str previous) {
        final StringBuilder result = new StringBuilder();
        Position position = previous.getPosition();
        final List<Str> list = assign.getValue();
        final int size = list.size();
        for (int i = 0; i < size; i++) {
            final Str item = list.get(i);
            result.append(item.toString());
            if (i != size - 1) {
                result.append(previous);
            }
            final Position pos = item.getPosition();
            if (pos != null && pos.getFile() != null) {
                position = pos;
            }
        }
        return new Str(position, result.toString());
    }
}
