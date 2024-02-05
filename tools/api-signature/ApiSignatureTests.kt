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
