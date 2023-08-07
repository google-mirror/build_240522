import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertThrows;
import static org.junit.Assert.assertTrue;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

import com.android.aconfig.test.FeatureFlags;
import com.android.aconfig.test.FeatureFlagsImpl;
import com.android.aconfig.test.Flags;

@RunWith(JUnit4.class)
public final class AconfigHostTest {
    @Test
    public void testThrowsExceptionIfFlagNotSet() {
        assertThrows(NullPointerException.class, () -> Flags.disabledRo());
        FeatureFlags featureFlags = new FeatureFlagsImpl();
        assertThrows(IllegalArgumentException.class, () -> featureFlags.disabledRo());
    }

    @Test
    public void testSetFlagInFeatureFlagsImpl() {
        FeatureFlags featureFlags = new FeatureFlagsImpl();
        featureFlags.setFlag(Flags.FLAG_ENABLED_RW, true);
        assertTrue(featureFlags.enabledRw());
        featureFlags.setFlag(Flags.FLAG_ENABLED_RW, false);
        assertFalse(featureFlags.enabledRw());

        //Set Flags
        assertThrows(NullPointerException.class, () -> Flags.enabledRw());
        Flags.setFeatureFlagsImpl(featureFlags);
        featureFlags.setFlag(Flags.FLAG_ENABLED_RW, true);
        assertTrue(Flags.enabledRw());
        Flags.unsetFeatureFlagsImpl();
    }

    @Test
    public void testSetFlagWithRandomName() {
        FeatureFlags featureFlags = new FeatureFlagsImpl();
        assertThrows(IllegalArgumentException.class,
            () -> featureFlags.setFlag("Randome_name", true));
    }

    @Test
    public void testResetFlagsInFeatureFlagsImpl() {
        FeatureFlags featureFlags = new FeatureFlagsImpl();
        featureFlags.setFlag(Flags.FLAG_ENABLED_RO, true);
        assertTrue(featureFlags.enabledRo());
        featureFlags.resetAll();
        assertThrows(IllegalArgumentException.class, () -> featureFlags.enabledRo());

        // Set value after reset
        featureFlags.setFlag(Flags.FLAG_ENABLED_RO, false);
        assertFalse(featureFlags.enabledRo());
    }

    @Test
    public void testFlagsSetFeatureFlagsImpl() {
        FeatureFlags featureFlags = new FeatureFlagsImpl();
        featureFlags.setFlag(Flags.FLAG_ENABLED_RW, true);
        assertThrows(NullPointerException.class, () -> Flags.enabledRw());
        Flags.setFeatureFlagsImpl(featureFlags);
        assertTrue(Flags.enabledRw());
        Flags.unsetFeatureFlagsImpl();
    }

    @Test
    public void testFlagsUnsetFeatureFlagsImpl() {
        FeatureFlags featureFlags = new FeatureFlagsImpl();
        featureFlags.setFlag(Flags.FLAG_ENABLED_RW, true);
        assertThrows(NullPointerException.class, () -> Flags.enabledRw());
        Flags.setFeatureFlagsImpl(featureFlags);
        assertTrue(Flags.enabledRw());

        Flags.unsetFeatureFlagsImpl();
        assertThrows(NullPointerException.class, () -> Flags.enabledRw());
    }
}
