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

import java.io.PrintStream;
import java.util.List;

public class MakeWriter {
    public static final int FLAG_WRITE_HEADER = 1;
    public static final int FLAG_WRITE_ANNOTATIONS = 1 << 1;

    private final boolean mWriteHeader;
    private final boolean mWriteAnnotations;

    public static void write(PrintStream out, GenericConfig config, int flags) {
        (new MakeWriter(flags)).write(out, config);
    }

    private MakeWriter(int flags) {
        mWriteHeader = (flags & FLAG_WRITE_HEADER) != 0;
        mWriteAnnotations = (flags & FLAG_WRITE_ANNOTATIONS) != 0;
    }

    private void write(PrintStream out, GenericConfig config) {
        for (GenericConfig.ConfigFile file: config.getFiles().values()) {
            out.println("---------------------------------------------------------");
            out.println("FILE: " + file.getFilename());
            out.println("---------------------------------------------------------");
            writeFile(out, config, file);
            out.println();
        }
    }

    private void writeFile(PrintStream out, GenericConfig config, GenericConfig.ConfigFile file) {
        if (mWriteHeader) {
            out.println("# This file is generated by the product_config tool");
        }
        for (GenericConfig.Statement statement: file.getStatements()) {
            if (statement instanceof GenericConfig.Assign) {
                writeAssign(out, config, (GenericConfig.Assign)statement);
            } else if (statement instanceof GenericConfig.Inherit) {
                writeInherit(out, (GenericConfig.Inherit)statement);
            } else {
                throw new RuntimeException("Unexpected Statement: " + statement);
            }
        }
    }

    private void writeAssign(PrintStream out, GenericConfig config,
            GenericConfig.Assign statement) {
        final List<Str> values = statement.getValue();
        final int size = values.size();
        final String varName = statement.getName();
        Position pos = null;
        if (size == 0) {
            return;
        } else if (size == 1) {
            // Plain :=
            final Str value = values.get(0);
            out.print(varName + " := " + value);
            pos = value.getPosition();
        } else if (size == 2 && values.get(0).toString().length() == 0) {
            // Plain +=
            final Str value = values.get(1);
            out.print(varName + " += " + value);
            pos = value.getPosition();
        } else {
            // Write it out the long way
            out.print(varName + " := " + values.get(0));
            for (int i = 1; i < size; i++) {
                out.print("$(" + varName + ") " + values.get(i));
                pos = values.get(i).getPosition();
            }
        }
        if (mWriteAnnotations) {
            out.print("  # " + config.getVarType(varName) + " " + pos);
        }
        out.println();
    }

    private void writeInherit(PrintStream out, GenericConfig.Inherit statement) {
        final Str filename = statement.getFilename();
        out.print("$(call inherit-product " + filename + ")");
        if (mWriteAnnotations) {
            out.print("  # " + filename.getPosition());
        }
        out.println();
    }
}
