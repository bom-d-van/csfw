// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package directory

import (
	"fmt"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

type (
	// SourceCurrencyAll used in Path: `system/currency/installed`,
	SourceCurrencyAll struct {
		mc config.ModelConstructor
	}
)

var _ config.FieldSourceModeller = (*SourceCurrencyAll)(nil)

// NewSourceCurrencyAll creates a new option for all currencies. If one argument of
// the ModelConstructor has been provided you may skip the calling of Construct().
func NewSourceCurrencyAll(mc ...config.ModelConstructor) *SourceCurrencyAll {
	sca := &SourceCurrencyAll{}
	if len(mc) == 1 {
		if err := sca.Construct(mc[0]); err != nil {
			log.Error("SourceCurrencyAll=NewSourceCurrencyAll", "err", err)
		}
	}
	return sca
}

// Construct sets the necessary options
func (sca *SourceCurrencyAll) Construct(mc config.ModelConstructor) error {
	if mc.ConfigReader == nil {
		return errgo.New("ConfigReader is required")
	}
	if mc.Scope == nil {
		return errgo.New("Scope is required")
	}
	sca.mc = mc
	return nil
}
func (sca *SourceCurrencyAll) Options() config.ValueLabelSlice {
	// Magento\Framework\Locale\Resolver
	// grep locale from general/locale/code scope::store for the current store ID
	// the store locale greps the currencies from http://php.net/manual/en/class.resourcebundle.php
	// in the correct language
	storeLocale := sca.mc.ConfigReader.GetString(config.Path(PathDefaultLocale), config.ScopeStore(sca.mc.Scope))

	fmt.Printf("\nstoreLocale: %s\n", storeLocale)

	return nil
}
