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

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.ProgramResult

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

class CheckCommand : CliktCommand() {
  override fun run() {
    println("hello world")
    throw ProgramResult(0)
  }
}

fun main(args: Array<String>) = CheckCommand().main(args)
