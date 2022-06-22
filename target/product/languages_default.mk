#
# Copyright (C) 2009 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Note: To make the language visible in the Language settings in Android,
# it must be defined both in this file (in the format xx_YY or xx_yyyy_ZZ)
# and in platform/frameworks/base/core/res/res/values-xx-rYY/,
# where xx is the language code and YY is the country code (all caps),
# or in platform/frameworks/base/core/res/res/values-b+xx+yyyy+ZZ/,
# or in platform/frameworks/base/core/res/res/values-b+xx+ZZ/,
# where xx is the language code, yyyy is the variant code (optional, first letter capitalized),
# and ZZ is the country code (all caps).

# This is a build configuration that just contains a list of languages, with
# en_US (English (United States)) set as the default language.
PRODUCT_LOCALES := \
        en_US \         # English (United States)
        af_ZA \         # Afrikaans
        am_ET \         # Amharic
        ar_EG \         # Arabic (Egypt)
        ar_XB \         # Arabic (Pseudo-Bidi)
        as_IN \         # Assamese
        az_AZ \         # Azerbaijani
        be_BY \         # Belarusian
        bg_BG \         # Bulgarian
        bn_BD \         # Bangla
        bs_BA \         # Bosnian
        ca_ES \         # Catalan
        cs_CZ \         # Czech
        da_DK \         # Danish
        de_AT \         # German (Austria)
        de_CH \         # German (Switzerland)
        de_DE \         # German (Germany)
        de_LI \         # German (Liechtenstein)
        el_GR \         # Greek
        en_AU \         # English (Australia)
        en_CA \         # English (Canada)
        en_GB \         # English (United Kingdom)
        en_IE \         # English (Ireland)
        en_IN \         # English (India)
        en_NZ \         # English (New Zealand)
        en_SG \         # English (Singapore)
        en_XA \         # English (Pseudo-Accents)
        en_ZA \         # English (South Africa)
        es_ES \         # Spanish (Spain)
        es_MX \         # Spanish (Mexico)
        es_US \         # Spanish (United States)
        et_EE \         # Estonian
        eu_ES \         # Basque
        fa_IR \         # Persian
        fi_FI \         # Finnish
        fr_BE \         # French (Belgium)
        fr_CA \         # French (Canada)
        fr_CH \         # French (Switzerland)
        fr_FR \         # French (France)
        gl_ES \         # Galician
        gu_IN \         # Gujarati
        hi_IN \         # Hindi
        hr_HR \         # Croatian
        hu_HU \         # Hungarian
        hy_AM \         # Armenian
        in_ID \         # Indonesian
        is_IS \         # Icelandic
        it_CH \         # Italian (Switzerland)
        it_IT \         # Italian (Italy)
        iw_IL \         # Hebrew
        ja_JP \         # Japanese
        ka_GE \         # Georgian
        kk_KZ \         # Kazakh
        km_KH \         # Khmer
        kn_IN \         # Kannada
        ko_KR \         # Korean
        ky_KG \         # Kyrgyz
        lo_LA \         # Lao
        lt_LT \         # Lithuanian
        lv_LV \         # Latvian
        mk_MK \         # Macedonian
        ml_IN \         # Malayalam
        mn_MN \         # Mongolian
        mr_IN \         # Marathi
        ms_MY \         # Malay
        my_MM \         # Burmese
        nb_NO \         # Norwegian Bokmål
        ne_NP \         # Nepali
	nl_BE \         # Dutch (Belgium)
        nl_NL \         # Dutch (Netherlands)
        or_IN \         # Odia
        pa_IN \         # Punjabi
        pl_PL \         # Polish
        pt_BR \         # Portuguese (Brazil)
        pt_PT \         # Portuguese (Portugal)
        ro_RO \         # Romanian
        ru_RU \         # Russian
        si_LK \         # Sinhala
        sk_SK \         # Slovak
        sl_SI \         # Slovenian
        sq_AL \         # Albanian
        sr_Latn_RS \    # Serbian (Latin, Serbia)
        sr_RS \         # Serbian (Serbia)
        sv_SE \         # Swedish
        sw_TZ \         # Swahili
        ta_IN \         # Tamil
        te_IN \         # Telugu
        th_TH \         # Thai
        tl_PH \         # Filipino/Tagalog
        tr_TR \         # Turkish
        uk_UA \         # Ukrainian
        ur_PK \         # Urdu
        uz_UZ \         # Uzbek
        vi_VN \         # Vietnamese
        zh_CN \         # Chinese (Simplified, China)
        zh_HK \         # Chinese (Traditional, Hong Kong)
        zh_TW \         # Chinese (Traditional, Taiwan)
        zu_ZA \         # Zulu
