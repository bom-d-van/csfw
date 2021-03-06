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

package store

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	// Storager implements the requirements to get new websites, groups and store views.
	// This interface is used in the StoreManager
	Storager interface {
		// Website creates a new Website pointer from an ID or code including all of its
		// groups and all related stores. It panics when the integrity is incorrect.
		// If ID and code are available then the non-empty code has precedence.
		Website(config.ScopeIDer) (*Website, error)
		// Websites creates a slice containing all pointers to Websites with its associated
		// groups and stores. It panics when the integrity is incorrect.
		Websites() (WebsiteSlice, error)
		// Group creates a new Group which contains all related stores and its website.
		// Only the argument ID can be used to get a specific Group.
		Group(config.ScopeIDer) (*Group, error)
		// Groups creates a slice containing all pointers to Groups with its associated
		// stores and websites. It panics when the integrity is incorrect.
		Groups() (GroupSlice, error)
		// Store creates a new Store containing its group and its website.
		// If ID and code are available then the non-empty code has precedence.
		Store(config.ScopeIDer) (*Store, error)
		// Stores creates a new store slice. Can return an error when the website or
		// the group cannot be found.
		Stores() (StoreSlice, error)
		// DefaultStoreView traverses through the websites to find the default website and gets
		// the default group which has the default store id assigned to. Only one website can be the default one.
		DefaultStoreView() (*Store, error)
		// ReInit reloads the websites, groups and stores from the database.
		ReInit(dbr.SessionRunner, ...csdb.DbrSelectCb) error
	}

	// Storage contains a mutex and the raw slices from the database. @todo maybe make private?
	Storage struct {
		cr       config.Reader
		mu       sync.RWMutex
		websites TableWebsiteSlice
		groups   TableGroupSlice
		stores   TableStoreSlice
	}

	// StorageOption option func for NewStorage()
	StorageOption func(*Storage)
)

// check if interface has been implemented
var _ Storager = (*Storage)(nil)

// SetStorageWebsites adds the TableWebsiteSlice to the Storage. By default, the slice is nil.
func SetStorageWebsites(tws ...*TableWebsite) StorageOption {
	return func(s *Storage) { s.websites = TableWebsiteSlice(tws) }
}

// SetStorageGroups adds the TableGroupSlice to the Storage. By default, the slice is nil.
func SetStorageGroups(tgs ...*TableGroup) StorageOption {
	return func(s *Storage) { s.groups = TableGroupSlice(tgs) }
}

// SetStorageStores adds the TableStoreSlice to the Storage. By default, the slice is nil.
func SetStorageStores(tss ...*TableStore) StorageOption {
	return func(s *Storage) { s.stores = TableStoreSlice(tss) }
}

// SetStorageConfig sets the configuration Reader. Optional.
// Default reader is config.DefaultManager
func SetStorageConfig(cr config.Reader) StorageOption {
	return func(s *Storage) { s.cr = cr }
}

// NewStorage creates a new storage object from three slice types. All three arguments can be nil
// but then you call ReInit()
func NewStorage(opts ...StorageOption) *Storage {
	s := &Storage{
		cr: config.DefaultManager,
		mu: sync.RWMutex{},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	return s
}

// NewStorageOption sames as NewStorage() but returns a function to be used in NewManager()
func NewStorageOption(opts ...StorageOption) ManagerOption {
	return func(m *Manager) { m.storage = NewStorage(opts...) }
}

// website returns a TableWebsite by using either id or code to find it. If id and code are
// available then the non-empty code has precedence.
func (st *Storage) website(r config.ScopeIDer) (*TableWebsite, error) {
	if r == nil {
		return nil, ErrWebsiteNotFound
	}
	if c, ok := r.(config.ScopeCoder); ok && c.ScopeCode() != "" {
		return st.websites.FindByCode(c.ScopeCode())
	}
	return st.websites.FindByID(r.ScopeID())
}

// Website creates a new Website according to the interface definition.
func (st *Storage) Website(r config.ScopeIDer) (*Website, error) {
	w, err := st.website(r)
	if err != nil {
		return nil, err
	}
	return NewWebsite(w).SetGroupsStores(st.groups, st.stores), nil
}

// Websites creates a slice of Website pointers according to the interface definition.
func (st *Storage) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(st.websites), len(st.websites))
	for i, w := range st.websites {
		websites[i] = NewWebsite(w).SetGroupsStores(st.groups, st.stores)
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument. If no argument or more than
// one has been supplied it returns an error.
func (st *Storage) group(r config.ScopeIDer) (*TableGroup, error) {
	if r == nil {
		return nil, ErrGroupNotFound
	}
	return st.groups.FindByID(r.ScopeID())
}

// Group creates a new Group which contains all related stores and its website according to the
// interface definition.
func (st *Storage) Group(id config.ScopeIDer) (*Group, error) {
	g, err := st.group(id)
	if err != nil {
		return nil, err
	}

	w, err := st.website(config.ScopeID(g.WebsiteID))
	if err != nil {
		return nil, err
	}
	return NewGroup(g, SetGroupWebsite(w), SetGroupConfig(st.cr)).SetStores(st.stores, nil), nil
}

// Groups creates a new group slice containing its website all related stores.
// May panic when a website pointer is nil.
func (st *Storage) Groups() (GroupSlice, error) {
	groups := make(GroupSlice, len(st.groups), len(st.groups))
	for i, g := range st.groups {
		w, err := st.website(config.ScopeID(g.WebsiteID))
		if err != nil {
			return nil, errgo.Mask(err)
		}
		groups[i] = NewGroup(g, SetGroupConfig(st.cr), SetGroupWebsite(w)).SetStores(st.stores, nil)
	}
	return groups, nil
}

// store returns a TableStore by an id or code.
// The non-empty code has precedence if available.
func (st *Storage) store(r config.ScopeIDer) (*TableStore, error) {
	if r == nil {
		return nil, ErrStoreNotFound
	}
	if c, ok := r.(config.ScopeCoder); ok && c.ScopeCode() != "" {
		return st.stores.FindByCode(c.ScopeCode())
	}
	return st.stores.FindByID(r.ScopeID())
}

// Store creates a new Store which contains the the store, its group and website
// according to the interface definition.
func (st *Storage) Store(r config.ScopeIDer) (*Store, error) {
	s, err := st.store(r)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	w, err := st.website(config.ScopeID(s.WebsiteID))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	g, err := st.group(config.ScopeID(s.GroupID))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	ns := NewStore(s, w, g, SetStoreConfig(st.cr))
	ns.Website().SetGroupsStores(st.groups, st.stores)
	ns.Group().SetStores(st.stores, w)
	return ns, nil
}

// Stores creates a new store slice. Can return an error when the website or
// the group cannot be found.
func (st *Storage) Stores() (StoreSlice, error) {
	stores := make(StoreSlice, len(st.stores), len(st.stores))
	for i, s := range st.stores {
		var err error
		if stores[i], err = st.Store(config.ScopeID(s.StoreID)); err != nil {
			return nil, errgo.Mask(err)
		}
	}
	return stores, nil
}

// DefaultStoreView traverses through the websites to find the default website and gets
// the default group which has the default store id assigned to. Only one website can be the default one.
func (st *Storage) DefaultStoreView() (*Store, error) {
	for _, website := range st.websites {
		if website.IsDefault.Bool && website.IsDefault.Valid {
			g, err := st.group(config.ScopeID(website.DefaultGroupID))
			if err != nil {
				return nil, err
			}
			return st.Store(config.ScopeID(g.DefaultStoreID))
		}
	}
	return nil, ErrStoreNotFound
}

// ReInit reloads all websites, groups and stores concurrently from the database. If GOMAXPROCS
// is set to > 1 then in parallel. Returns an error with location or nil. If an error occurs
// then all internal slices will be reset.
func (st *Storage) ReInit(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	errc := make(chan error)
	defer close(errc)
	// not sure about those three go
	go func() {
		for i := range st.websites {
			st.websites[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.websites = nil
		_, err := st.websites.Load(dbrSess, cbs...)
		errc <- errgo.Mask(err)
	}()

	go func() {
		for i := range st.groups {
			st.groups[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.groups = nil
		_, err := st.groups.Load(dbrSess, cbs...)
		errc <- errgo.Mask(err)
	}()

	go func() {
		for i := range st.stores {
			st.stores[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.stores = nil
		_, err := st.stores.Load(dbrSess, cbs...)
		errc <- errgo.Mask(err)
	}()

	for i := 0; i < 3; i++ {
		if err := <-errc; err != nil {
			// in case of error clear all
			st.websites = nil
			st.groups = nil
			st.stores = nil
			return err
		}
	}
	return nil
}
