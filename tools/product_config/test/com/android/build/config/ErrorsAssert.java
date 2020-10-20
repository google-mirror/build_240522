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

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import java.nio.charset.StandardCharsets;

public class ErrorsAssert {
    public static void assertHasEntry(Errors.Category category, Errors errors) {
        StringBuilder found = new StringBuilder();
        for (Errors.Entry entry: errors.getEntries()) {
            if (entry.getCategory() == category) {
                return;
            }
            found.append(' ');
            found.append(entry.getCategory().getCode());
        }
        throw new AssertionError("No error category " + category.getCode() + " found."
                + " Found category codes were:" + found);
    }

    public static String getErrorMessages(Errors errors) {
        final ByteArrayOutputStream stream = new ByteArrayOutputStream();
        try {
            errors.printErrors(new PrintStream(stream, true, StandardCharsets.UTF_8.name()));
        } catch (UnsupportedEncodingException ex) {
            // utf-8 is always supported
        }
        return new String(stream.toByteArray(), StandardCharsets.UTF_8);
    }
}

