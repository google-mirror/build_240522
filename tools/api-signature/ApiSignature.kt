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
