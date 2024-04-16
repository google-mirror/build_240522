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
@file:JvmName("Main")

package com.android.checkflaggedapis

import com.android.tools.metalava.model.BaseItemVisitor
import com.android.tools.metalava.model.FieldItem
import com.android.tools.metalava.model.text.ApiFile
import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.ProgramResult
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.required
import com.github.ajalt.clikt.parameters.types.path
import java.io.InputStream

@JvmInline
value class Symbol(val name: String) {
  companion object {
    private val FORBIDDEN_CHARS = listOf('/', '#', '$')

    fun create(name: String): Symbol {
      var sanitized_name = name
      for (ch in FORBIDDEN_CHARS) {
        sanitized_name = sanitized_name.replace(ch, '.')
      }
      return Symbol(sanitized_name)
    }
  }

  init {
    require(!name.isEmpty()) { "empty string" }
    for (ch in FORBIDDEN_CHARS) {
      require(!name.contains(ch)) { "$name: contains $ch" }
    }
  }

  override fun toString(): String = name.toString()
}

@JvmInline
value class Flag(val name: String) {
  override fun toString(): String = name.toString()
}

class CheckCommand : CliktCommand() {
  private val api_signature_path by
      option("--api-signature")
          .path(mustExist = true, canBeDir = false, mustBeReadable = true)
          .required()

  override fun run() {
    @Suppress("UNUSED_VARIABLE")
    val flagged_symbols =
        api_signature_path.toFile().inputStream().use { inputStream ->
          parseApiSignature(api_signature_path.toString(), inputStream)
        }
    throw ProgramResult(0)
  }
}

private fun parseApiSignature(path: String, input: InputStream): Set<Pair<Symbol, Flag>> {
  // TODO(334870672): add support for classes and metods
  val output = mutableSetOf<Pair<Symbol, Flag>>()
  val visitor =
      object : BaseItemVisitor() {
        override fun visitField(field: FieldItem) {
          val flag =
              field.modifiers
                  .findAnnotation("android.annotation.FlaggedApi")
                  ?.findAttribute("value")
                  ?.value
                  ?.value() as? String
          if (flag != null) {
            val symbol = Symbol.create(field.baselineElementId())
            output.add(Pair(symbol, Flag(flag)))
          }
        }
      }
  val codebase = ApiFile.parseApi(path, input)
  codebase.accept(visitor)
  return output
}

fun main(args: Array<String>) = CheckCommand().main(args)
