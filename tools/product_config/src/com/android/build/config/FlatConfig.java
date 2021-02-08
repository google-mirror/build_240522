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

import java.io.PrintStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;

/**
 * Flattened configuration -- set of variables after all assignments and inherits have
 * been executed.
 */
public class FlatConfig extends ConfigBase {

    private final TreeMap<String, Value> mValues = new TreeMap();

    /**
     * Class to hold the two types of variables we support, strings and lists of strings.
     */
    public static class Value {
        private final VarType mVarType;
        private final Str mStr;
        private final ArrayList<Str> mList;

        public Value(VarType varType, Str str) {
            mVarType = varType;
            mStr = str;
            mList = null;
        }

        public Value(List<Str> list) {
            mVarType = VarType.LIST;
            mStr = null;
            mList = new ArrayList(list);
        }

        public VarType getVarType() {
            return mVarType;
        }

        public Str getStr() {
            return mStr;
        }

        public List<Str> getList() {
            return mList;
        }

        public String debugString() {
            final StringBuilder str = new StringBuilder("Value(type=");
            str.append(mVarType.toString());
            str.append(" mStr=");
            if (mStr == null) {
                str.append("null");
            } else {
                str.append("\"");
                str.append(mStr.toString());
                str.append("\" (");
                str.append(" (");
                str.append(mStr.getPosition().toString());
                str.append(")");
            }
            str.append(" mList=");
            if (mList == null) {
                str.append("null");
            } else {
                str.append("[");
                for (Str s: mList) {
                    str.append("\"");
                    str.append(s.toString());
                    str.append("\" (");
                    str.append(s.getPosition().toString());
                    str.append(")");
                }
                str.append(" ]");
            }
            str.append(")");
            return str.toString();
        }
    }

    public TreeMap<String, Value> getValues() {
        return mValues;
    }
}
