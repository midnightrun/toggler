/*
	Package specs



	Summary

	This package implements generic CRUD definitions that commonly appears upon interacting with external resources.



	Minimum Requirement from resource point of view

	In order to make this package work, you have to implement the TestMinimumRequirementsWithExampleEntities specification.
	Most of the other resource specs specification depends on the resource specs mentioned in the min requirement specification.
	Keep in mind, that you have no guarantee on your resource content during test execution, because some specification
	may alter the content of the resource (db), or delete from it.
	If you need specific data in the resource you want to test with, you must ensure in the test execution that
	such context is correctly provisioned, and after test execution, cleaned up.
	If you use such data-set in a external resource that needs to be kept intact,
	I advise you to use separate environments for test execution and manual testing.



	Requirement from Business Entities

	This package depends on a fact that there is a string field ID in a business entity struct,
	or at least a tag `ext:"ID"`. This allows the package to create specifications that assumes,
	that the ID field links the EntityType structure to an external resource object.
	The resource specs package doesn't care about the content of the ID string field,
	and don't have assumptions other than the existence of the field ID on a struct


*/
package specs
