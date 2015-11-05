// Copyright (C) 2015 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//

#include <assert.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>

#include <string>

#include "dex_file-inl.h"
#include "mem_map.h"

using art::DexFile;

namespace {

void DumpField(const DexFile& dex_file, uint32_t idx, uint32_t flags,
               art::EncodedStaticFieldValueIterator* static_field_values) {
  const DexFile::FieldId& field_id = dex_file.GetFieldId(idx);
  const char* name = dex_file.StringDataByIdx(field_id.name_idx_);
  const char* type_desc = dex_file.StringByTypeIdx(field_id.type_idx_);
  const char* class_desc = dex_file.StringByTypeIdx(field_id.class_idx_);
  printf("%s (flag=%04x, type='%s', class='%s')",
         name, flags, type_desc, class_desc);

  if (static_field_values) {
    const jvalue& v = static_field_values->GetJavaValue();
    switch (static_field_values->GetValueType()) {
      case art::EncodedStaticFieldValueIterator::kByte:
        printf(" = %" PRIu8, v.b);
        break;
      case art::EncodedStaticFieldValueIterator::kShort:
        printf(" = %" PRId16, v.s);
        break;
      case art::EncodedStaticFieldValueIterator::kChar:
        printf(" = %" PRIu16, v.c);
        break;
      case art::EncodedStaticFieldValueIterator::kInt:
        printf(" = %" PRId32, v.i);
        break;
      case art::EncodedStaticFieldValueIterator::kLong:
        printf(" = %" PRId64, v.j);
        break;
      case art::EncodedStaticFieldValueIterator::kFloat:
        printf(" = %f", v.f);
        break;
      case art::EncodedStaticFieldValueIterator::kDouble:
        printf(" = %f", v.d);
        break;
      case art::EncodedStaticFieldValueIterator::kString:
        // TODO
        printf(" = %s", "STRING");
        break;
      case art::EncodedStaticFieldValueIterator::kType:
        // TODO
        printf(" = %s", "TYPE");
        break;
      case art::EncodedStaticFieldValueIterator::kNull:
        printf(" = null");
        break;
      case art::EncodedStaticFieldValueIterator::kBoolean:
        printf(" = %s", v.z ? "true" : "false");
        break;

      case art::EncodedStaticFieldValueIterator::kField:
      case art::EncodedStaticFieldValueIterator::kMethod:
      case art::EncodedStaticFieldValueIterator::kEnum:
      case art::EncodedStaticFieldValueIterator::kArray:
      case art::EncodedStaticFieldValueIterator::kAnnotation:
      default:
        assert(false);
    }
  }

  puts("");
}

void DumpMethod(const DexFile& dex_file, uint32_t idx, uint32_t flags) {
  const DexFile::MethodId& method_id = dex_file.GetMethodId(idx);
  const char* name = dex_file.StringDataByIdx(method_id.name_idx_);
  const art::Signature signature = dex_file.GetMethodSignature(method_id);
  char* type_desc = strdup(signature.ToString().c_str());
  const char* class_desc = dex_file.StringByTypeIdx(method_id.class_idx_);
  printf("%s (flag=%04x, type'=%s', class='%s')\n",
         name, flags, type_desc, class_desc);
}

void DumpClassData(const DexFile& dex_file,
                   const DexFile::ClassDef& class_def) {
  const uint8_t* encoded_data = dex_file.GetClassData(class_def);
  if (!encoded_data)
    return;
  art::ClassDataItemIterator class_data(dex_file, encoded_data);

  art::EncodedStaticFieldValueIterator static_field_values(dex_file, class_def);
  for (; class_data.HasNextStaticField(); class_data.Next()) {
    printf(" Static field: ");
    assert(static_field_values.HasNext());
    DumpField(dex_file, class_data.GetMemberIndex(),
              class_data.GetRawMemberAccessFlags(),
              &static_field_values);
    static_field_values.Next();
  }
  assert(!static_field_values.HasNext());
  for (; class_data.HasNextInstanceField(); class_data.Next()) {
    printf(" Instance field: ");
    DumpField(dex_file, class_data.GetMemberIndex(),
              class_data.GetRawMemberAccessFlags(), nullptr);
  }
  for (; class_data.HasNextDirectMethod(); class_data.Next()) {
    printf(" Direct method: ");
    DumpMethod(dex_file,
               class_data.GetMemberIndex(),
               class_data.GetRawMemberAccessFlags());
  }
  for (; class_data.HasNextVirtualMethod(); class_data.Next()) {
    printf(" Virtual method: ");
    DumpMethod(dex_file,
               class_data.GetMemberIndex(),
               class_data.GetRawMemberAccessFlags());
  }
}

void DumpClass(const DexFile& dex_file, const DexFile::ClassDef& class_def) {
  const char* class_desc = dex_file.StringByTypeIdx(class_def.class_idx_);
  const char* superclass_desc = nullptr;
  if (class_def.superclass_idx_ != DexFile::kDexNoIndex16) {
    dex_file.StringByTypeIdx(class_def.superclass_idx_);
  }

  printf("\n");
  printf("Class: '%s'\n", class_desc);
  printf("Access flags: %0x04x\n", class_def.access_flags_);
  if (superclass_desc) {
    printf("Super class: '%s'\n", superclass_desc);
  }

  const DexFile::TypeList* interfaces = dex_file.GetInterfacesList(class_def);
  if (interfaces != nullptr) {
    for (uint32_t i = 0; i < interfaces->Size(); i++) {
      const char* interface_desc =
          dex_file.StringByTypeIdx(interfaces->GetTypeItem(i).type_idx_);
      printf("Interface: '%s'\n", interface_desc);
    }
  }

  DumpClassData(dex_file, class_def);
}

void ProcessDexFile(const DexFile& dex_file) {
  printf("DEX version: %s\n", dex_file.GetHeader().magic_ + 4);
  for (uint32_t i = 0; i < dex_file.GetHeader().class_defs_size_; i++) {
    const DexFile::ClassDef& class_def = dex_file.GetClassDef(i);
    DumpClass(dex_file, class_def);
  }
}

}  // namespace

int main(int argc, char* argv[]) {
  if (argc != 2) {
    fprintf(stderr, "Usage: %s input.dex\n", argv[0]);
    exit(1);
  }

  art::InitLogging(argv);
  art::MemMap::Init();

  const char* filename = argv[1];
  std::string error_msg;
  std::vector<std::unique_ptr<const DexFile>> dex_files;
  if (!DexFile::Open(filename, filename, &error_msg, &dex_files)) {
    fprintf(stderr, "%s\n", error_msg.c_str());
    exit(1);
  }

  for (size_t i = 0; i < dex_files.size(); i++) {
    const DexFile& dex_file = *dex_files[i];
    printf("Filename: %s\n", filename);
    ProcessDexFile(dex_file);
  }
}
