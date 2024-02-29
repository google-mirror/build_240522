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

import com.beust.jcommander.JCommander
import com.beust.jcommander.Parameter
import com.beust.jcommander.Parameters
import com.beust.jcommander.converters.FileConverter
import java.io.File

const val EXTRACT_API_FLAGS = "extract-api-flags"

@Parameters(
  commandNames = arrayOf(EXTRACT_API_FLAGS),
  commandDescription =
    "Extract the aconfig flags that are used together with @FlaggedApi from an API signature file"
)
class FooCommand {
  @Parameter(
    names = arrayOf("--api-file"),
    description = "Path to API signature file",
    required = true,
    converter = FileConverter::class
  )
  lateinit var path: File

  fun call() {
    val contents = path.readText(Charsets.UTF_8)
    val flags = find_flagged_api_flags(path.absolutePath, contents)
    println(flags.sorted().joinToString("\n"))
  }
}

fun main(args: Array<String>) {
  val fooCommand = FooCommand()
  val jc = JCommander.newBuilder().addCommand(fooCommand).build()
  jc.parse(*args)
  when (jc.parsedCommand) {
    EXTRACT_API_FLAGS -> fooCommand.call()
    else -> jc.usage()
  }
}
