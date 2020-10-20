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

import org.junit.Assert;
import org.junit.Test;

import java.util.HashSet;
import java.util.List;

public class ErrorsTest {
    /**
     * Test that errors can be recorded and retrieved.
     */
    @Test
    public void testAdding() {
        Errors errors = new Errors();

        errors.add(errors.ERROR_COMMAND_LINE, new Position("a", 12), "Errrororrrr");

        Assert.assertTrue(errors.hadWarningOrError());
        Assert.assertTrue(errors.hadError());

        List<Errors.Entry> entries = errors.getEntries();
        Assert.assertEquals(1, entries.size());

        Errors.Entry entry = entries.get(0);
        Assert.assertEquals(errors.ERROR_COMMAND_LINE, entry.getCategory());
        Assert.assertEquals("a", entry.getPosition().getFile());
        Assert.assertEquals(12, entry.getPosition().getLine());
        Assert.assertEquals("Errrororrrr", entry.getMessage());
    }

    /**
     * Test that not adding an error doesn't record errors.
     */
    @Test
    public void testNoError() {
        Errors errors = new Errors();

        Assert.assertFalse(errors.hadWarningOrError());
        Assert.assertFalse(errors.hadError());
    }

    /**
     * Test that not adding a warning doesn't record errors.
     */
    @Test
    public void testWarning() {
        Errors errors = new Errors();

        errors.add(errors.REPLACE_ME, "Waaaaarninggggg");
        
        Assert.assertTrue(errors.hadWarningOrError());
        Assert.assertFalse(errors.hadError());
    }
}
