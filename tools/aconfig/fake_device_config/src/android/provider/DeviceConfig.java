/*
 * Copyright (C) 2023 The Android Open Source Project
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

package android.provider;

/*
 * This class allows generated aconfig code to compile independently of the framework.
 */
public class DeviceConfig {
	private DeviceConfig() {
	}

	public static boolean getBoolean(String ns, String name, boolean def) {
		return false;
	}

	public static class Properties {
		public boolean getBoolean(String name, boolean def) {
			return false;
		}
	}
}
