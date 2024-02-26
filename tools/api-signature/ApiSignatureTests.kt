/*
 * Copyright (C) 2024 The Android Open Source Project
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
package android.apisignature

import org.junit.Assert.assertEquals
import org.junit.Test

class ApiSignatureTests {
  @Test
  fun `extract API flags`() {
    val contents =
      """// Signature format: 2.0
         package android {
           public final class Foo {
             ctor public Foo();
             field @FlaggedApi("android.foo") public static final String FOO = "android.FOO";
             method @FlaggedApi("android.bar") public void foo();
           }
         }
      """
    val flags = find_flagged_api_flags(":in-memory:", contents)
    assertEquals(setOf("android.bar", "android.foo"), flags)
  }
}
