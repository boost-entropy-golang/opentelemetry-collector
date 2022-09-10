// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal // import "go.opentelemetry.io/collector/pdata/internal/cmd/pdatagen/internal"

import (
	"os"
	"strings"
)

const accessorSliceTemplate = `// ${fieldName} returns the ${originFieldName} associated with this ${structName}.
func (ms ${structName}) ${fieldName}() ${packageName}${returnType} {
	return ${packageName}${returnType}(internal.New${returnType}(&ms.getOrig().${originFieldName}))
}`

const accessorsSliceTestTemplate = `func Test${structName}_${fieldName}(t *testing.T) {
	ms := New${structName}()
	assert.Equal(t, ${packageName}New${returnType}(), ms.${fieldName}())
	internal.FillTest${returnType}(internal.${returnType}(ms.${fieldName}()))
	assert.Equal(t, ${packageName}${returnType}(internal.GenerateTest${returnType}()), ms.${fieldName}())
}`

const accessorsMessageValueTemplate = `// ${fieldName} returns the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) ${fieldName}() ${packageName}${returnType} {
	return ${packageName}${returnType}(internal.New${returnType}(&ms.getOrig().${originFieldName}))
}`

const accessorsMessageValueTestTemplate = `func Test${structName}_${fieldName}(t *testing.T) {
	ms := New${structName}()
	internal.FillTest${returnType}(internal.${returnType}(ms.${fieldName}()))
	assert.Equal(t, ${packageName}${returnType}(internal.GenerateTest${returnType}()), ms.${fieldName}())
}`

const accessorsPrimitiveTemplate = `// ${fieldName} returns the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) ${fieldName}() ${packageName}${returnType} {
	return ms.getOrig().${originFieldName}
}

// Set${fieldName} replaces the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) Set${fieldName}(v ${returnType}) {
	ms.getOrig().${originFieldName} = v
}`

const accessorsPrimitiveSliceTemplate = `// ${fieldName} returns the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) ${fieldName}() ${packageName}${returnType} {
	return ${packageName}${returnType}(internal.New${returnType}(&ms.getOrig().${originFieldName}))
}

// Set${fieldName} replaces the ${lowerFieldName} associated with this ${structName}.
// Deprecated: [0.60.0] Use ${fieldName}().FromRaw() instead
func (ms ${structName}) Set${fieldName}(v ${packageName}${returnType}) {
	ms.getOrig().${originFieldName} = *internal.GetOrig${returnType}(internal.${returnType}(v))
}`

const oneOfTypeAccessorHeaderTemplate = `// ${originFieldName}Type returns the type of the ${lowerOriginFieldName} for this ${structName}.
// Calling this function on zero-initialized ${structName} will cause a panic.
func (ms ${structName}) ${originFieldName}Type() ${typeName} {
	switch ms.getOrig().${originFieldName}.(type) {`

const oneOfTypeAccessorHeaderTestTemplate = `func Test${structName}_${originFieldName}Type(t *testing.T) {
	tv := New${structName}()
	assert.Equal(t, ${typeName}None, tv.${originFieldName}Type())
}
`

const accessorsOneOfMessageTemplate = `// ${fieldName} returns the ${lowerFieldName} associated with this ${structName}.
//
// Calling this function when ${originOneOfFieldName}Type() != ${typeName} returns an invalid 
// zero-initialized instance of ${returnType}. Note that using such ${returnType} instance can cause panic.
//
// Calling this function on zero-initialized ${structName} will cause a panic.
func (ms ${structName}) ${fieldName}() ${returnType} {
	v, ok := ms.getOrig().Get${originOneOfFieldName}().(*${originStructType})
	if !ok {
		return ${returnType}{}
	}
	return new${returnType}(v.${originFieldName})
}

// SetEmpty${fieldName} sets an empty ${lowerFieldName} to this ${structName}.
//
// After this, ${originOneOfFieldName}Type() function will return ${typeName}".
//
// Calling this function on zero-initialized ${structName} will cause a panic.
func (ms ${structName}) SetEmpty${fieldName}() ${returnType} {
	val := &${originFieldPackageName}.${originFieldName}{}
	ms.getOrig().${originOneOfFieldName} = &${originStructType}{${originFieldName}: val}
	return new${returnType}(val)
}`

const accessorsOneOfMessageTestTemplate = `func Test${structName}_${fieldName}(t *testing.T) {
	ms := New${structName}()
	internal.FillTest${returnType}(internal.${returnType}(ms.SetEmpty${fieldName}()))
	assert.Equal(t, ${typeName}, ms.${originOneOfFieldName}Type())
	assert.Equal(t, ${returnType}(internal.GenerateTest${returnType}()), ms.${fieldName}())
}

func Test${structName}_CopyTo_${fieldName}(t *testing.T) {
	ms := New${structName}()
	internal.FillTest${returnType}(internal.${returnType}(ms.SetEmpty${fieldName}()))
	dest := New${structName}()
	ms.CopyTo(dest)
	assert.Equal(t, ms, dest)
}`

const copyToValueOneOfMessageTemplate = `	case ${typeName}:
		ms.${fieldName}().CopyTo(dest.SetEmpty${fieldName}())`

const accessorsOneOfPrimitiveTemplate = `// ${fieldName} returns the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) ${fieldName}() ${returnType} {
	return ms.getOrig().Get${originFieldName}()
}

// Set${fieldName} replaces the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) Set${fieldName}(v ${returnType}) {
	ms.getOrig().${originOneOfFieldName} = &${originStructType}{
		${originFieldName}: v,
	}
}`

const accessorsOneOfPrimitiveTestTemplate = `func Test${structName}_${fieldName}(t *testing.T) {
	ms := New${structName}()
	assert.Equal(t, ${defaultVal}, ms.${fieldName}())
	ms.Set${fieldName}(${testValue})
	assert.Equal(t, ${testValue}, ms.${fieldName}())
	assert.Equal(t, ${typeName}, ms.${originOneOfFieldName}Type())
}`

const accessorsPrimitiveTestTemplate = `func Test${structName}_${fieldName}(t *testing.T) {
	ms := New${structName}()
	assert.Equal(t, ${defaultVal}, ms.${fieldName}())
	ms.Set${fieldName}(${testValue})
	assert.Equal(t, ${testValue}, ms.${fieldName}())
}`

const accessorsPrimitiveTypedTemplate = `// ${fieldName} returns the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) ${fieldName}() ${packageName}${returnType} {
	return ${packageName}${returnType}(ms.getOrig().${originFieldName})
}

// Set${fieldName} replaces the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) Set${fieldName}(v ${packageName}${returnType}) {
	ms.getOrig().${originFieldName} = ${rawType}(v)
}`

const accessorsPrimitiveTypedTestTemplate = `func Test${structName}_${fieldName}(t *testing.T) {
	ms := New${structName}()
	assert.Equal(t, ${packageName}${returnType}(${defaultVal}), ms.${fieldName}())
	testVal${fieldName} := ${packageName}${returnType}(${testValue})
	ms.Set${fieldName}(testVal${fieldName})
	assert.Equal(t, testVal${fieldName}, ms.${fieldName}())
}`

const accessorsPrimitiveSliceTestTemplate = `func Test${structName}_${fieldName}(t *testing.T) {
	ms := New${structName}()
	assert.Equal(t, ${defaultVal}, ms.${fieldName}().AsRaw())
	ms.${fieldName}().FromRaw(${testValue})
	assert.Equal(t, ${testValue}, ms.${fieldName}().AsRaw())
}`

const accessorsOptionalPrimitiveValueTemplate = `// ${fieldName} returns the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) ${fieldName}() ${returnType} {
	return ms.getOrig().Get${fieldName}()
}
// Has${fieldName} returns true if the ${structName} contains a
// ${fieldName} value, false otherwise.
func (ms ${structName}) Has${fieldName}() bool {
	return ms.getOrig().${fieldName}_ != nil
}
// Set${fieldName} replaces the ${lowerFieldName} associated with this ${structName}.
func (ms ${structName}) Set${fieldName}(v ${returnType}) {
	ms.getOrig().${fieldName}_ = &${originStructType}{${fieldName}: v}
}`

type baseField interface {
	generateAccessors(ms baseStruct, sb *strings.Builder)

	generateAccessorsTest(ms baseStruct, sb *strings.Builder)

	generateSetWithTestValue(sb *strings.Builder)

	generateCopyToValue(ms baseStruct, sb *strings.Builder)
}

type sliceField struct {
	fieldName       string
	originFieldName string
	returnSlice     baseSlice
}

func (sf *sliceField) generateAccessors(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorSliceTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return sf.fieldName
		case "packageName":
			if sf.returnSlice.getPackageName() != ms.getPackageName() {
				return sf.returnSlice.getPackageName() + "."
			}
			return ""
		case "returnType":
			return sf.returnSlice.getName()
		case "originFieldName":
			return sf.originFieldName
		default:
			panic(name)
		}
	}))
}

func (sf *sliceField) generateAccessorsTest(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsSliceTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return sf.fieldName
		case "packageName":
			if sf.returnSlice.getPackageName() != ms.getPackageName() {
				return sf.returnSlice.getPackageName() + "."
			}
			return ""
		case "returnType":
			return sf.returnSlice.getName()
		default:
			panic(name)
		}
	}))
}

func (sf *sliceField) generateSetWithTestValue(sb *strings.Builder) {
	sb.WriteString("\tFillTest" + sf.returnSlice.getName() + "(New" + sf.returnSlice.getName() + "(&tv.orig." + sf.originFieldName + "))")
}

func (sf *sliceField) generateCopyToValue(_ baseStruct, sb *strings.Builder) {
	sb.WriteString("\tms." + sf.fieldName + "().CopyTo(dest." + sf.fieldName + "())")
}

var _ baseField = (*sliceField)(nil)

type messageValueField struct {
	fieldName       string
	originFieldName string
	returnMessage   baseStruct
}

func (mf *messageValueField) generateAccessors(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsMessageValueTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return mf.fieldName
		case "lowerFieldName":
			return strings.ToLower(mf.fieldName)
		case "packageName":
			if mf.returnMessage.getPackageName() != ms.getPackageName() {
				return mf.returnMessage.getPackageName() + "."
			}
			return ""
		case "returnType":
			return mf.returnMessage.getName()
		case "originFieldName":
			return mf.originFieldName
		default:
			panic(name)
		}
	}))
}

func (mf *messageValueField) generateAccessorsTest(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsMessageValueTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return mf.fieldName
		case "returnType":
			return mf.returnMessage.getName()
		case "packageName":
			if mf.returnMessage.getPackageName() != ms.getPackageName() {
				return mf.returnMessage.getPackageName() + "."
			}
			return ""
		default:
			panic(name)
		}
	}))
}

func (mf *messageValueField) generateSetWithTestValue(sb *strings.Builder) {
	sb.WriteString("\tFillTest" + mf.returnMessage.getName() + "(New" + mf.returnMessage.getName() + "(&tv.orig." + mf.originFieldName + "))")
}

func (mf *messageValueField) generateCopyToValue(_ baseStruct, sb *strings.Builder) {
	sb.WriteString("\tms." + mf.fieldName + "().CopyTo(dest." + mf.fieldName + "())")
}

var _ baseField = (*messageValueField)(nil)

type primitiveField struct {
	fieldName       string
	originFieldName string
	returnType      string
	defaultVal      string
	testVal         string
}

func (pf *primitiveField) generateAccessors(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsPrimitiveTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "packageName":
			return ""
		case "fieldName":
			return pf.fieldName
		case "lowerFieldName":
			return strings.ToLower(pf.fieldName)
		case "returnType":
			return pf.returnType
		case "originFieldName":
			return pf.originFieldName
		default:
			panic(name)
		}
	}))
}

func (pf *primitiveField) generateAccessorsTest(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsPrimitiveTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "packageName":
			return ""
		case "defaultVal":
			return pf.defaultVal
		case "fieldName":
			return pf.fieldName
		case "testValue":
			return pf.testVal
		default:
			panic(name)
		}
	}))
}

func (pf *primitiveField) generateSetWithTestValue(sb *strings.Builder) {
	sb.WriteString("\ttv.orig." + pf.originFieldName + " = " + pf.testVal)
}

func (pf *primitiveField) generateCopyToValue(_ baseStruct, sb *strings.Builder) {
	sb.WriteString("\tdest.Set" + pf.fieldName + "(ms." + pf.fieldName + "())")
}

var _ baseField = (*primitiveField)(nil)

type primitiveType struct {
	structName  string
	packageName string
	rawType     string
	defaultVal  string
	testVal     string
}

// Types that has defined a custom type (e.g. "type Timestamp uint64")
type primitiveTypedField struct {
	fieldName       string
	originFieldName string
	returnType      *primitiveType
}

func (ptf *primitiveTypedField) generateAccessors(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsPrimitiveTypedTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return ptf.fieldName
		case "lowerFieldName":
			return strings.ToLower(ptf.fieldName)
		case "returnType":
			return ptf.returnType.structName
		case "packageName":
			if ptf.returnType.packageName != ms.getPackageName() {
				return ptf.returnType.packageName + "."
			}
			return ""
		case "rawType":
			return ptf.returnType.rawType
		case "originFieldName":
			return ptf.originFieldName
		default:
			panic(name)
		}
	}))
}

func (ptf *primitiveTypedField) generateAccessorsTest(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsPrimitiveTypedTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "defaultVal":
			return ptf.returnType.defaultVal
		case "packageName":
			if ptf.returnType.packageName != ms.getPackageName() {
				return ptf.returnType.packageName + "."
			}
			return ""
		case "returnType":
			return ptf.returnType.structName
		case "fieldName":
			return ptf.fieldName
		case "testValue":
			return ptf.returnType.testVal
		default:
			panic(name)
		}
	}))
}

func (ptf *primitiveTypedField) generateSetWithTestValue(sb *strings.Builder) {
	sb.WriteString("\ttv.orig." + ptf.originFieldName + " = " + ptf.returnType.testVal)
}

func (ptf *primitiveTypedField) generateCopyToValue(_ baseStruct, sb *strings.Builder) {
	sb.WriteString("\tdest.Set" + ptf.fieldName + "(ms." + ptf.fieldName + "())")
}

var _ baseField = (*primitiveTypedField)(nil)

// primitiveSliceField is used to generate fields for slice of primitive types
type primitiveSliceField struct {
	fieldName         string
	originFieldName   string
	returnPackageName string
	returnType        string
	defaultVal        string
	rawType           string
	testVal           string
}

func (psf *primitiveSliceField) generateAccessors(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsPrimitiveSliceTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return psf.fieldName
		case "lowerFieldName":
			return strings.ToLower(psf.fieldName)
		case "returnType":
			return psf.returnType
		case "packageName":
			if psf.returnPackageName != ms.getPackageName() {
				return psf.returnPackageName + "."
			}
			return ""
		case "originFieldName":
			return psf.originFieldName
		default:
			panic(name)
		}
	}))
}

func (psf *primitiveSliceField) generateAccessorsTest(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsPrimitiveSliceTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "packageName":
			if psf.returnPackageName != ms.getPackageName() {
				return psf.returnPackageName + "."
			}
			return ""
		case "returnType":
			return psf.returnType
		case "defaultVal":
			return psf.defaultVal
		case "fieldName":
			return psf.fieldName
		case "testValue":
			return psf.testVal
		default:
			panic(name)
		}
	}))
}

func (psf *primitiveSliceField) generateSetWithTestValue(sb *strings.Builder) {
	sb.WriteString("\ttv.orig." + psf.originFieldName + " = " + psf.testVal)
}

func (psf *primitiveSliceField) generateCopyToValue(ms baseStruct, sb *strings.Builder) {
	sb.WriteString("\tms." + psf.fieldName + "().CopyTo(dest." + psf.fieldName + "())")
}

var _ baseField = (*primitiveSliceField)(nil)

type oneOfField struct {
	originTypePrefix string
	originFieldName  string
	typeName         string
	testValueIdx     int
	values           []oneOfValue
}

func (of *oneOfField) generateAccessors(ms baseStruct, sb *strings.Builder) {
	of.generateTypeAccessors(ms, sb)
	sb.WriteString("\n")
	for _, v := range of.values {
		v.generateAccessors(ms, of, sb)
		sb.WriteString("\n")
	}
}

func (of *oneOfField) generateTypeAccessors(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(oneOfTypeAccessorHeaderTemplate, func(name string) string {
		switch name {
		case "lowerOriginFieldName":
			return strings.ToLower(of.originFieldName)
		case "originFieldName":
			return of.originFieldName
		case "structName":
			return ms.getName()
		case "typeName":
			return of.typeName
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
	for _, v := range of.values {
		v.generateTypeSwitchCase(of, sb)
	}
	sb.WriteString("\t}\n")
	sb.WriteString("\treturn " + of.typeName + "None\n")
	sb.WriteString("}\n")
}

func (of *oneOfField) generateAccessorsTest(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(oneOfTypeAccessorHeaderTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "originFieldName":
			return of.originFieldName
		case "typeName":
			return of.typeName
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
	for _, v := range of.values {
		v.generateTests(ms, of, sb)
		sb.WriteString("\n")
	}
}

func (of *oneOfField) generateSetWithTestValue(sb *strings.Builder) {
	of.values[of.testValueIdx].generateSetWithTestValue(of, sb)
}

func (of *oneOfField) generateCopyToValue(_ baseStruct, sb *strings.Builder) {
	sb.WriteString("\tswitch ms." + of.originFieldName + "Type() {\n")
	for _, v := range of.values {
		v.generateCopyToValue(of, sb)
	}
	sb.WriteString("\t}\n")
}

var _ baseField = (*oneOfField)(nil)

type oneOfValue interface {
	getFieldType() string
	generateAccessors(ms baseStruct, of *oneOfField, sb *strings.Builder)
	generateTests(ms baseStruct, of *oneOfField, sb *strings.Builder)
	generateSetWithTestValue(of *oneOfField, sb *strings.Builder)
	generateCopyToValue(of *oneOfField, sb *strings.Builder)
	generateTypeSwitchCase(of *oneOfField, sb *strings.Builder)
}

type oneOfPrimitiveValue struct {
	fieldName       string
	fieldType       string
	defaultVal      string
	testVal         string
	returnType      string
	originFieldName string
}

func (opv *oneOfPrimitiveValue) getFieldType() string {
	return opv.fieldType
}

func (opv *oneOfPrimitiveValue) generateAccessors(ms baseStruct, of *oneOfField, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsOneOfPrimitiveTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return opv.fieldName
		case "lowerFieldName":
			return strings.ToLower(opv.fieldName)
		case "returnType":
			return opv.returnType
		case "originFieldName":
			return opv.originFieldName
		case "originOneOfFieldName":
			return of.originFieldName
		case "originStructType":
			return of.originTypePrefix + opv.originFieldName
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
}

func (opv *oneOfPrimitiveValue) generateTests(ms baseStruct, of *oneOfField, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsOneOfPrimitiveTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "defaultVal":
			return opv.defaultVal
		case "packageName":
			return ""
		case "fieldName":
			return opv.fieldName
		case "testValue":
			return opv.testVal
		case "originOneOfFieldName":
			return of.originFieldName
		case "typeName":
			return of.typeName + opv.fieldType
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
}

func (opv *oneOfPrimitiveValue) generateSetWithTestValue(of *oneOfField, sb *strings.Builder) {
	sb.WriteString("\ttv.orig." + of.originFieldName + " = &" + of.originTypePrefix + opv.originFieldName + "{" + opv.originFieldName + ":" + opv.testVal + "}")
}

func (opv *oneOfPrimitiveValue) generateCopyToValue(of *oneOfField, sb *strings.Builder) {
	sb.WriteString("\tcase " + of.typeName + opv.fieldType + ":\n")
	sb.WriteString("\tdest.Set" + opv.fieldName + "(ms." + opv.fieldName + "())\n")
}

func (opv *oneOfPrimitiveValue) generateTypeSwitchCase(of *oneOfField, sb *strings.Builder) {
	sb.WriteString("\tcase *" + of.originTypePrefix + opv.originFieldName + ":\n")
	sb.WriteString("\t\treturn " + of.typeName + opv.fieldType + "\n")
}

var _ oneOfValue = (*oneOfPrimitiveValue)(nil)

type oneOfMessageValue struct {
	fieldName              string
	originFieldName        string
	originFieldPackageName string
	returnMessage          *messageValueStruct
}

func (omv *oneOfMessageValue) getFieldType() string {
	return omv.fieldName
}

func (omv *oneOfMessageValue) generateAccessors(ms baseStruct, of *oneOfField, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsOneOfMessageTemplate, func(name string) string {
		switch name {
		case "fieldName":
			return omv.fieldName
		case "lowerFieldName":
			return strings.ToLower(omv.fieldName)
		case "originFieldName":
			return omv.originFieldName
		case "originOneOfFieldName":
			return of.originFieldName
		case "originFieldPackageName":
			return omv.originFieldPackageName
		case "originStructType":
			return of.originTypePrefix + omv.originFieldName
		case "returnType":
			return omv.returnMessage.structName
		case "structName":
			return ms.getName()
		case "typeName":
			return of.typeName + omv.returnMessage.structName
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
}

func (omv *oneOfMessageValue) generateTests(ms baseStruct, of *oneOfField, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsOneOfMessageTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return omv.fieldName
		case "returnType":
			return omv.returnMessage.structName
		case "originOneOfFieldName":
			return of.originFieldName
		case "typeName":
			return of.typeName + omv.returnMessage.structName
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
}

func (omv *oneOfMessageValue) generateSetWithTestValue(of *oneOfField, sb *strings.Builder) {
	sb.WriteString("\ttv.orig." + of.originFieldName + " = &" + of.originTypePrefix + omv.originFieldName + "{" + omv.originFieldName + ": &" + omv.originFieldPackageName + "." + omv.originFieldName + "{}}\n")
	sb.WriteString("\tFillTest" + omv.returnMessage.structName + "(New" + omv.fieldName + "(tv.orig.Get" + omv.originFieldName + "()))")
}

func (omv *oneOfMessageValue) generateCopyToValue(of *oneOfField, sb *strings.Builder) {
	sb.WriteString(os.Expand(copyToValueOneOfMessageTemplate, func(name string) string {
		switch name {
		case "fieldName":
			return omv.fieldName
		case "originOneOfFieldName":
			return of.originFieldName
		case "typeName":
			return of.typeName + omv.fieldName
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
}

func (omv *oneOfMessageValue) generateTypeSwitchCase(of *oneOfField, sb *strings.Builder) {
	sb.WriteString("\tcase *" + of.originTypePrefix + omv.originFieldName + ":\n")
	sb.WriteString("\t\treturn " + of.typeName + omv.fieldName + "\n")
}

var _ oneOfValue = (*oneOfMessageValue)(nil)

type optionalPrimitiveValue struct {
	fieldName        string
	defaultVal       string
	testVal          string
	returnType       string
	originFieldName  string
	originTypePrefix string
}

func (opv *optionalPrimitiveValue) generateAccessors(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsOptionalPrimitiveValueTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "fieldName":
			return opv.fieldName
		case "lowerFieldName":
			return strings.ToLower(opv.fieldName)
		case "returnType":
			return opv.returnType
		case "originFieldName":
			return opv.originFieldName
		case "originStructType":
			return opv.originTypePrefix + opv.originFieldName
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
}

func (opv *optionalPrimitiveValue) generateAccessorsTest(ms baseStruct, sb *strings.Builder) {
	sb.WriteString(os.Expand(accessorsPrimitiveTestTemplate, func(name string) string {
		switch name {
		case "structName":
			return ms.getName()
		case "packageName":
			return ""
		case "defaultVal":
			return opv.defaultVal
		case "fieldName":
			return opv.fieldName
		case "testValue":
			return opv.testVal
		default:
			panic(name)
		}
	}))
	sb.WriteString("\n")
}

func (opv *optionalPrimitiveValue) generateSetWithTestValue(sb *strings.Builder) {
	sb.WriteString("\ttv.orig." + opv.originFieldName + "_ = &" + opv.originTypePrefix + opv.originFieldName + "{" + opv.originFieldName + ":" + opv.testVal + "}")
}

func (opv *optionalPrimitiveValue) generateCopyToValue(_ baseStruct, sb *strings.Builder) {
	sb.WriteString("if ms.Has" + opv.fieldName + "(){\n")
	sb.WriteString("\tdest.Set" + opv.fieldName + "(ms." + opv.fieldName + "())\n")
	sb.WriteString("}\n")
}

var _ baseField = (*optionalPrimitiveValue)(nil)
