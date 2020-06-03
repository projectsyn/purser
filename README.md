# Purser: a tool to check access permissions and preconditions on clouds. 

Purser takes the credentials for a cloud providers API and checks several
aspects against a list of expectations. In its initial form, it is written to
test the requirements the OpenShift 4 installer has for a Google Compute
Platform project.

See [CHANGELOG.md](/CHANGELOG.md) for changelogs of each release version of
Purser.

## Usage

Purser needs to now where to finde the service account keys. To this by
exporting the variable `GOOGLE_APPLICATION_CREDENTIALS` pointing to the JSON
file holding this data. 

    purser gcp <project id> --domain <domain>

The `<project>` id can be passed in its numerical or human readable form. You
can find the project ID in the JSON file holding the service account key.

The domain is optional. If set, validation will fail if the given domain is not
a public managed zone within the validated project. The domain has to be given
as FQDN including the dot at the end.
