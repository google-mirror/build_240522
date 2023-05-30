import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertTrue;
import static test.Flags.disabled_ro;
import static test.Flags.disabled_rw;
import static test.Flags.enabled_ro;
import static test.Flags.enabled_rw;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

@RunWith(JUnit4.class)
public final class AconfigTest {
  @Test
  public void testDisabledReadOnlyFlag() {
    assertFalse(disabled_ro());
  }

  @Test
  public void testEnabledReadOnlyFlag() {
    assertTrue(enabled_ro());
  }

  @Test
  public void testDisabledReadWriteFlag() {
    assertFalse(disabled_rw());
  }

  @Test
  public void testEnabledReadWriteFlag() {
    assertTrue(enabled_rw());
  }
}
