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

import com.android.tools.metalava.model.BaseItemVisitor
import com.android.tools.metalava.model.Item
import com.android.tools.metalava.model.text.ApiFile
import java.io.File

fun main(args: Array<String>) {
  for (path in args) {
    val contents = File(path).readText(Charsets.UTF_8)
    val flags = find_flagged_api_flags(path, contents)
    println(flags.sorted().joinToString("\n"))
  }
}

// Throws: com.android.tools.metalava.model.text.ApiParseException
internal fun find_flagged_api_flags(path: String, contents: String): Set<String> {
  val flagSet: MutableSet<String> = mutableSetOf()
  val visitor =
    object : BaseItemVisitor() {
      override fun visitItem(item: Item) {
        val value =
          item.modifiers
            .findAnnotation("android.annotation.FlaggedApi")
            ?.findAttribute("value")
            ?.value
            ?.value() as? String
        if (value != null) {
          flagSet.add(value)
        }
      }
    }

  val codebase = ApiFile.parseApi(path, contents)
  codebase.accept(visitor)
  return flagSet
}
