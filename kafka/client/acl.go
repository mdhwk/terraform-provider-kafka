package client

import (
	"context"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kmsg"
)

type ACLResult struct {
	Principal  *string
	Host       *string
	Type       kmsg.ACLResourceType
	Name       *string
	Pattern    kadm.ACLPattern
	Operation  kadm.ACLOperation
	Permission kmsg.ACLPermissionType
	Err        error
}

func CreateACLs(ctx context.Context, c *kadm.Client, acl *kadm.ACLBuilder) ([]ACLResult, error) {
	if err := acl.ValidateCreate(); err != nil {
		return nil, err
	}

	response, err := c.CreateACLs(ctx, acl)
	if err != nil {
		return nil, err
	}

	res := make([]ACLResult, len(response))
	for i, r := range response {
		res[i] = ACLResult{
			Principal:  &r.Principal,
			Host:       &r.Host,
			Type:       r.Type,
			Name:       &r.Name,
			Pattern:    r.Pattern,
			Operation:  r.Operation,
			Permission: r.Permission,
			Err:        r.Err,
		}
	}

	return res, nil
}

func DeleteACLs(ctx context.Context, c *kadm.Client, acl *kadm.ACLBuilder) (res []ACLResult, err error) {
	if err = acl.ValidateDelete(); err != nil {
		return nil, err
	}

	response, err := c.DeleteACLs(ctx, acl)
	if err != nil {
		return nil, err
	}

	res = make([]ACLResult, len(response))
	for i, r := range response {
		res[i] = ACLResult{
			Principal:  r.Principal,
			Host:       r.Host,
			Type:       r.Type,
			Name:       r.Name,
			Pattern:    r.Pattern,
			Operation:  r.Operation,
			Permission: r.Permission,
			Err:        r.Err,
		}
	}

	return res, nil
}

func ReadACLs(ctx context.Context, c *kadm.Client, acl *kadm.ACLBuilder) (res []ACLResult, err error) {
	if err = acl.ValidateDescribe(); err != nil {
		return nil, err
	}

	response, err := c.DescribeACLs(ctx, acl)
	if err != nil {
		return nil, err
	}

	res = make([]ACLResult, len(response))
	for i, r := range response {
		res[i] = ACLResult{
			Principal:  r.Principal,
			Host:       r.Host,
			Type:       r.Type,
			Name:       r.Name,
			Pattern:    r.Pattern,
			Operation:  r.Operation,
			Permission: r.Permission,
			Err:        r.Err,
		}
	}

	return res, nil
}
