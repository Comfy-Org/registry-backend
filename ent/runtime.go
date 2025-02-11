// Code generated by ent, DO NOT EDIT.

package ent

import (
	"registry-backend/ent/ciworkflowresult"
	"registry-backend/ent/comfynode"
	"registry-backend/ent/gitcommit"
	"registry-backend/ent/node"
	"registry-backend/ent/nodereview"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/personalaccesstoken"
	"registry-backend/ent/publisher"
	"registry-backend/ent/schema"
	"registry-backend/ent/storagefile"
	"registry-backend/ent/user"
	"time"

	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	ciworkflowresultMixin := schema.CIWorkflowResult{}.Mixin()
	ciworkflowresultMixinFields0 := ciworkflowresultMixin[0].Fields()
	_ = ciworkflowresultMixinFields0
	ciworkflowresultFields := schema.CIWorkflowResult{}.Fields()
	_ = ciworkflowresultFields
	// ciworkflowresultDescCreateTime is the schema descriptor for create_time field.
	ciworkflowresultDescCreateTime := ciworkflowresultMixinFields0[0].Descriptor()
	// ciworkflowresult.DefaultCreateTime holds the default value on creation for the create_time field.
	ciworkflowresult.DefaultCreateTime = ciworkflowresultDescCreateTime.Default.(func() time.Time)
	// ciworkflowresultDescUpdateTime is the schema descriptor for update_time field.
	ciworkflowresultDescUpdateTime := ciworkflowresultMixinFields0[1].Descriptor()
	// ciworkflowresult.DefaultUpdateTime holds the default value on creation for the update_time field.
	ciworkflowresult.DefaultUpdateTime = ciworkflowresultDescUpdateTime.Default.(func() time.Time)
	// ciworkflowresult.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	ciworkflowresult.UpdateDefaultUpdateTime = ciworkflowresultDescUpdateTime.UpdateDefault.(func() time.Time)
	// ciworkflowresultDescStatus is the schema descriptor for status field.
	ciworkflowresultDescStatus := ciworkflowresultFields[5].Descriptor()
	// ciworkflowresult.DefaultStatus holds the default value on creation for the status field.
	ciworkflowresult.DefaultStatus = schema.WorkflowRunStatusType(ciworkflowresultDescStatus.Default.(string))
	// ciworkflowresultDescID is the schema descriptor for id field.
	ciworkflowresultDescID := ciworkflowresultFields[0].Descriptor()
	// ciworkflowresult.DefaultID holds the default value on creation for the id field.
	ciworkflowresult.DefaultID = ciworkflowresultDescID.Default.(func() uuid.UUID)
	comfynodeMixin := schema.ComfyNode{}.Mixin()
	comfynodeMixinFields0 := comfynodeMixin[0].Fields()
	_ = comfynodeMixinFields0
	comfynodeFields := schema.ComfyNode{}.Fields()
	_ = comfynodeFields
	// comfynodeDescCreateTime is the schema descriptor for create_time field.
	comfynodeDescCreateTime := comfynodeMixinFields0[0].Descriptor()
	// comfynode.DefaultCreateTime holds the default value on creation for the create_time field.
	comfynode.DefaultCreateTime = comfynodeDescCreateTime.Default.(func() time.Time)
	// comfynodeDescUpdateTime is the schema descriptor for update_time field.
	comfynodeDescUpdateTime := comfynodeMixinFields0[1].Descriptor()
	// comfynode.DefaultUpdateTime holds the default value on creation for the update_time field.
	comfynode.DefaultUpdateTime = comfynodeDescUpdateTime.Default.(func() time.Time)
	// comfynode.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	comfynode.UpdateDefaultUpdateTime = comfynodeDescUpdateTime.UpdateDefault.(func() time.Time)
	// comfynodeDescDeprecated is the schema descriptor for deprecated field.
	comfynodeDescDeprecated := comfynodeFields[6].Descriptor()
	// comfynode.DefaultDeprecated holds the default value on creation for the deprecated field.
	comfynode.DefaultDeprecated = comfynodeDescDeprecated.Default.(bool)
	// comfynodeDescExperimental is the schema descriptor for experimental field.
	comfynodeDescExperimental := comfynodeFields[7].Descriptor()
	// comfynode.DefaultExperimental holds the default value on creation for the experimental field.
	comfynode.DefaultExperimental = comfynodeDescExperimental.Default.(bool)
	// comfynodeDescOutputIsList is the schema descriptor for output_is_list field.
	comfynodeDescOutputIsList := comfynodeFields[8].Descriptor()
	// comfynode.DefaultOutputIsList holds the default value on creation for the output_is_list field.
	comfynode.DefaultOutputIsList = comfynodeDescOutputIsList.Default.([]bool)
	// comfynodeDescID is the schema descriptor for id field.
	comfynodeDescID := comfynodeFields[0].Descriptor()
	// comfynode.DefaultID holds the default value on creation for the id field.
	comfynode.DefaultID = comfynodeDescID.Default.(func() uuid.UUID)
	gitcommitMixin := schema.GitCommit{}.Mixin()
	gitcommitMixinFields0 := gitcommitMixin[0].Fields()
	_ = gitcommitMixinFields0
	gitcommitFields := schema.GitCommit{}.Fields()
	_ = gitcommitFields
	// gitcommitDescCreateTime is the schema descriptor for create_time field.
	gitcommitDescCreateTime := gitcommitMixinFields0[0].Descriptor()
	// gitcommit.DefaultCreateTime holds the default value on creation for the create_time field.
	gitcommit.DefaultCreateTime = gitcommitDescCreateTime.Default.(func() time.Time)
	// gitcommitDescUpdateTime is the schema descriptor for update_time field.
	gitcommitDescUpdateTime := gitcommitMixinFields0[1].Descriptor()
	// gitcommit.DefaultUpdateTime holds the default value on creation for the update_time field.
	gitcommit.DefaultUpdateTime = gitcommitDescUpdateTime.Default.(func() time.Time)
	// gitcommit.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	gitcommit.UpdateDefaultUpdateTime = gitcommitDescUpdateTime.UpdateDefault.(func() time.Time)
	// gitcommitDescID is the schema descriptor for id field.
	gitcommitDescID := gitcommitFields[0].Descriptor()
	// gitcommit.DefaultID holds the default value on creation for the id field.
	gitcommit.DefaultID = gitcommitDescID.Default.(func() uuid.UUID)
	nodeMixin := schema.Node{}.Mixin()
	nodeMixinFields0 := nodeMixin[0].Fields()
	_ = nodeMixinFields0
	nodeFields := schema.Node{}.Fields()
	_ = nodeFields
	// nodeDescCreateTime is the schema descriptor for create_time field.
	nodeDescCreateTime := nodeMixinFields0[0].Descriptor()
	// node.DefaultCreateTime holds the default value on creation for the create_time field.
	node.DefaultCreateTime = nodeDescCreateTime.Default.(func() time.Time)
	// nodeDescUpdateTime is the schema descriptor for update_time field.
	nodeDescUpdateTime := nodeMixinFields0[1].Descriptor()
	// node.DefaultUpdateTime holds the default value on creation for the update_time field.
	node.DefaultUpdateTime = nodeDescUpdateTime.Default.(func() time.Time)
	// node.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	node.UpdateDefaultUpdateTime = nodeDescUpdateTime.UpdateDefault.(func() time.Time)
	// nodeDescTags is the schema descriptor for tags field.
	nodeDescTags := nodeFields[10].Descriptor()
	// node.DefaultTags holds the default value on creation for the tags field.
	node.DefaultTags = nodeDescTags.Default.([]string)
	// nodeDescTotalInstall is the schema descriptor for total_install field.
	nodeDescTotalInstall := nodeFields[11].Descriptor()
	// node.DefaultTotalInstall holds the default value on creation for the total_install field.
	node.DefaultTotalInstall = nodeDescTotalInstall.Default.(int64)
	// nodeDescTotalStar is the schema descriptor for total_star field.
	nodeDescTotalStar := nodeFields[12].Descriptor()
	// node.DefaultTotalStar holds the default value on creation for the total_star field.
	node.DefaultTotalStar = nodeDescTotalStar.Default.(int64)
	// nodeDescTotalReview is the schema descriptor for total_review field.
	nodeDescTotalReview := nodeFields[13].Descriptor()
	// node.DefaultTotalReview holds the default value on creation for the total_review field.
	node.DefaultTotalReview = nodeDescTotalReview.Default.(int64)
	nodereviewFields := schema.NodeReview{}.Fields()
	_ = nodereviewFields
	// nodereviewDescStar is the schema descriptor for star field.
	nodereviewDescStar := nodereviewFields[3].Descriptor()
	// nodereview.DefaultStar holds the default value on creation for the star field.
	nodereview.DefaultStar = nodereviewDescStar.Default.(int)
	// nodereviewDescID is the schema descriptor for id field.
	nodereviewDescID := nodereviewFields[0].Descriptor()
	// nodereview.DefaultID holds the default value on creation for the id field.
	nodereview.DefaultID = nodereviewDescID.Default.(func() uuid.UUID)
	nodeversionMixin := schema.NodeVersion{}.Mixin()
	nodeversionMixinFields0 := nodeversionMixin[0].Fields()
	_ = nodeversionMixinFields0
	nodeversionFields := schema.NodeVersion{}.Fields()
	_ = nodeversionFields
	// nodeversionDescCreateTime is the schema descriptor for create_time field.
	nodeversionDescCreateTime := nodeversionMixinFields0[0].Descriptor()
	// nodeversion.DefaultCreateTime holds the default value on creation for the create_time field.
	nodeversion.DefaultCreateTime = nodeversionDescCreateTime.Default.(func() time.Time)
	// nodeversionDescUpdateTime is the schema descriptor for update_time field.
	nodeversionDescUpdateTime := nodeversionMixinFields0[1].Descriptor()
	// nodeversion.DefaultUpdateTime holds the default value on creation for the update_time field.
	nodeversion.DefaultUpdateTime = nodeversionDescUpdateTime.Default.(func() time.Time)
	// nodeversion.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	nodeversion.UpdateDefaultUpdateTime = nodeversionDescUpdateTime.UpdateDefault.(func() time.Time)
	// nodeversionDescDeprecated is the schema descriptor for deprecated field.
	nodeversionDescDeprecated := nodeversionFields[5].Descriptor()
	// nodeversion.DefaultDeprecated holds the default value on creation for the deprecated field.
	nodeversion.DefaultDeprecated = nodeversionDescDeprecated.Default.(bool)
	// nodeversionDescStatusReason is the schema descriptor for status_reason field.
	nodeversionDescStatusReason := nodeversionFields[7].Descriptor()
	// nodeversion.DefaultStatusReason holds the default value on creation for the status_reason field.
	nodeversion.DefaultStatusReason = nodeversionDescStatusReason.Default.(string)
	// nodeversionDescComfyNodeExtractStatus is the schema descriptor for comfy_node_extract_status field.
	nodeversionDescComfyNodeExtractStatus := nodeversionFields[8].Descriptor()
	// nodeversion.DefaultComfyNodeExtractStatus holds the default value on creation for the comfy_node_extract_status field.
	nodeversion.DefaultComfyNodeExtractStatus = schema.ComfyNodeExtractStatus(nodeversionDescComfyNodeExtractStatus.Default.(string))
	// nodeversionDescID is the schema descriptor for id field.
	nodeversionDescID := nodeversionFields[0].Descriptor()
	// nodeversion.DefaultID holds the default value on creation for the id field.
	nodeversion.DefaultID = nodeversionDescID.Default.(func() uuid.UUID)
	personalaccesstokenMixin := schema.PersonalAccessToken{}.Mixin()
	personalaccesstokenMixinFields0 := personalaccesstokenMixin[0].Fields()
	_ = personalaccesstokenMixinFields0
	personalaccesstokenFields := schema.PersonalAccessToken{}.Fields()
	_ = personalaccesstokenFields
	// personalaccesstokenDescCreateTime is the schema descriptor for create_time field.
	personalaccesstokenDescCreateTime := personalaccesstokenMixinFields0[0].Descriptor()
	// personalaccesstoken.DefaultCreateTime holds the default value on creation for the create_time field.
	personalaccesstoken.DefaultCreateTime = personalaccesstokenDescCreateTime.Default.(func() time.Time)
	// personalaccesstokenDescUpdateTime is the schema descriptor for update_time field.
	personalaccesstokenDescUpdateTime := personalaccesstokenMixinFields0[1].Descriptor()
	// personalaccesstoken.DefaultUpdateTime holds the default value on creation for the update_time field.
	personalaccesstoken.DefaultUpdateTime = personalaccesstokenDescUpdateTime.Default.(func() time.Time)
	// personalaccesstoken.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	personalaccesstoken.UpdateDefaultUpdateTime = personalaccesstokenDescUpdateTime.UpdateDefault.(func() time.Time)
	// personalaccesstokenDescID is the schema descriptor for id field.
	personalaccesstokenDescID := personalaccesstokenFields[0].Descriptor()
	// personalaccesstoken.DefaultID holds the default value on creation for the id field.
	personalaccesstoken.DefaultID = personalaccesstokenDescID.Default.(func() uuid.UUID)
	publisherMixin := schema.Publisher{}.Mixin()
	publisherMixinFields0 := publisherMixin[0].Fields()
	_ = publisherMixinFields0
	publisherFields := schema.Publisher{}.Fields()
	_ = publisherFields
	// publisherDescCreateTime is the schema descriptor for create_time field.
	publisherDescCreateTime := publisherMixinFields0[0].Descriptor()
	// publisher.DefaultCreateTime holds the default value on creation for the create_time field.
	publisher.DefaultCreateTime = publisherDescCreateTime.Default.(func() time.Time)
	// publisherDescUpdateTime is the schema descriptor for update_time field.
	publisherDescUpdateTime := publisherMixinFields0[1].Descriptor()
	// publisher.DefaultUpdateTime holds the default value on creation for the update_time field.
	publisher.DefaultUpdateTime = publisherDescUpdateTime.Default.(func() time.Time)
	// publisher.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	publisher.UpdateDefaultUpdateTime = publisherDescUpdateTime.UpdateDefault.(func() time.Time)
	storagefileMixin := schema.StorageFile{}.Mixin()
	storagefileMixinFields0 := storagefileMixin[0].Fields()
	_ = storagefileMixinFields0
	storagefileFields := schema.StorageFile{}.Fields()
	_ = storagefileFields
	// storagefileDescCreateTime is the schema descriptor for create_time field.
	storagefileDescCreateTime := storagefileMixinFields0[0].Descriptor()
	// storagefile.DefaultCreateTime holds the default value on creation for the create_time field.
	storagefile.DefaultCreateTime = storagefileDescCreateTime.Default.(func() time.Time)
	// storagefileDescUpdateTime is the schema descriptor for update_time field.
	storagefileDescUpdateTime := storagefileMixinFields0[1].Descriptor()
	// storagefile.DefaultUpdateTime holds the default value on creation for the update_time field.
	storagefile.DefaultUpdateTime = storagefileDescUpdateTime.Default.(func() time.Time)
	// storagefile.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	storagefile.UpdateDefaultUpdateTime = storagefileDescUpdateTime.UpdateDefault.(func() time.Time)
	// storagefileDescBucketName is the schema descriptor for bucket_name field.
	storagefileDescBucketName := storagefileFields[1].Descriptor()
	// storagefile.BucketNameValidator is a validator for the "bucket_name" field. It is called by the builders before save.
	storagefile.BucketNameValidator = storagefileDescBucketName.Validators[0].(func(string) error)
	// storagefileDescObjectName is the schema descriptor for object_name field.
	storagefileDescObjectName := storagefileFields[2].Descriptor()
	// storagefile.ObjectNameValidator is a validator for the "object_name" field. It is called by the builders before save.
	storagefile.ObjectNameValidator = storagefileDescObjectName.Validators[0].(func(string) error)
	// storagefileDescFilePath is the schema descriptor for file_path field.
	storagefileDescFilePath := storagefileFields[3].Descriptor()
	// storagefile.FilePathValidator is a validator for the "file_path" field. It is called by the builders before save.
	storagefile.FilePathValidator = storagefileDescFilePath.Validators[0].(func(string) error)
	// storagefileDescFileType is the schema descriptor for file_type field.
	storagefileDescFileType := storagefileFields[4].Descriptor()
	// storagefile.FileTypeValidator is a validator for the "file_type" field. It is called by the builders before save.
	storagefile.FileTypeValidator = storagefileDescFileType.Validators[0].(func(string) error)
	// storagefileDescID is the schema descriptor for id field.
	storagefileDescID := storagefileFields[0].Descriptor()
	// storagefile.DefaultID holds the default value on creation for the id field.
	storagefile.DefaultID = storagefileDescID.Default.(func() uuid.UUID)
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreateTime is the schema descriptor for create_time field.
	userDescCreateTime := userMixinFields0[0].Descriptor()
	// user.DefaultCreateTime holds the default value on creation for the create_time field.
	user.DefaultCreateTime = userDescCreateTime.Default.(func() time.Time)
	// userDescUpdateTime is the schema descriptor for update_time field.
	userDescUpdateTime := userMixinFields0[1].Descriptor()
	// user.DefaultUpdateTime holds the default value on creation for the update_time field.
	user.DefaultUpdateTime = userDescUpdateTime.Default.(func() time.Time)
	// user.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	user.UpdateDefaultUpdateTime = userDescUpdateTime.UpdateDefault.(func() time.Time)
	// userDescIsApproved is the schema descriptor for is_approved field.
	userDescIsApproved := userFields[3].Descriptor()
	// user.DefaultIsApproved holds the default value on creation for the is_approved field.
	user.DefaultIsApproved = userDescIsApproved.Default.(bool)
	// userDescIsAdmin is the schema descriptor for is_admin field.
	userDescIsAdmin := userFields[4].Descriptor()
	// user.DefaultIsAdmin holds the default value on creation for the is_admin field.
	user.DefaultIsAdmin = userDescIsAdmin.Default.(bool)
}
