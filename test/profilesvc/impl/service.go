package impl

import "sync"
import (
	profilesvcapi "personal/jeremyxu/sandbox/test/profilesvc/api"
	"context"
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]profilesvcapi.Profile
}

func NewInmemService() profilesvcapi.Service {
	return &inmemService{
		m: map[string]profilesvcapi.Profile{},
	}
}

func (s *inmemService) PostProfile(ctx context.Context, p profilesvcapi.Profile) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[p.ID]; ok {
		return profilesvcapi.ErrAlreadyExists // POST = create, don't overwrite
	}
	s.m[p.ID] = p
	return nil
}

func (s *inmemService) GetProfile(ctx context.Context, id string) (profilesvcapi.Profile, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return profilesvcapi.Profile{}, profilesvcapi.ErrNotFound
	}
	return p, nil
}

func (s *inmemService) PutProfile(ctx context.Context, id string, p profilesvcapi.Profile) error {
	if id != p.ID {
		return profilesvcapi.ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[id] = p // PUT = create or update
	return nil
}

func (s *inmemService) PatchProfile(ctx context.Context, id string, p profilesvcapi.Profile) error {
	if p.ID != "" && id != p.ID {
		return profilesvcapi.ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	existing, ok := s.m[id]
	if !ok {
		return profilesvcapi.ErrNotFound // PATCH = update existing, don't create
	}

	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the Profile definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if p.Name != "" {
		existing.Name = p.Name
	}
	if len(p.Addresses) > 0 {
		existing.Addresses = p.Addresses
	}
	s.m[id] = existing
	return nil
}

func (s *inmemService) DeleteProfile(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return profilesvcapi.ErrNotFound
	}
	delete(s.m, id)
	return nil
}

func (s *inmemService) GetAddresses(ctx context.Context, profileID string) ([]profilesvcapi.Address, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[profileID]
	if !ok {
		return []profilesvcapi.Address{}, profilesvcapi.ErrNotFound
	}
	return p.Addresses, nil
}

func (s *inmemService) GetAddress(ctx context.Context, profileID string, addressID string) (profilesvcapi.Address, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[profileID]
	if !ok {
		return profilesvcapi.Address{}, profilesvcapi.ErrNotFound
	}
	for _, address := range p.Addresses {
		if address.ID == addressID {
			return address, nil
		}
	}
	return profilesvcapi.Address{}, profilesvcapi.ErrNotFound
}

func (s *inmemService) PostAddress(ctx context.Context, profileID string, a profilesvcapi.Address) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[profileID]
	if !ok {
		return profilesvcapi.ErrNotFound
	}
	for _, address := range p.Addresses {
		if address.ID == a.ID {
			return profilesvcapi.ErrAlreadyExists
		}
	}
	p.Addresses = append(p.Addresses, a)
	s.m[profileID] = p
	return nil
}

func (s *inmemService) DeleteAddress(ctx context.Context, profileID string, addressID string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[profileID]
	if !ok {
		return profilesvcapi.ErrNotFound
	}
	newAddresses := make([]profilesvcapi.Address, 0, len(p.Addresses))
	for _, address := range p.Addresses {
		if address.ID == addressID {
			continue // delete
		}
		newAddresses = append(newAddresses, address)
	}
	if len(newAddresses) == len(p.Addresses) {
		return profilesvcapi.ErrNotFound
	}
	p.Addresses = newAddresses
	s.m[profileID] = p
	return nil
}
