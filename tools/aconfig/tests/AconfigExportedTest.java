import static com.android.aconfig.test.Flags.FLAG_DISABLED_RW_EXPORTED;
import static org.junit.Assert.assertEquals;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

import com.android.aconfig.test.FakeFeatureFlagsImpl;
import com.android.aconfig.test.FeatureFlags;

@RunWith(JUnit4.class)
public final class AconfigTest {
    @Test
    public void testExportedFlag() {
        assertEquals("com.android.aconfig.test.disabled_rw_exported", FLAG_DISABLED_RW_EXPORTED);
        assertFalse(disabledRo());
    }
}
