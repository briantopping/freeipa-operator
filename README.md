# FreeIPA Operator

An integrated Identity and Authentication solution for Linux/UNIX networked environments. IPA stands for "Identity, Policy, Audit".
Comparatively, it provides a solution similar to Microsoft Active Directory for *NIX environments.

At it's core, FreeIPA uses the 389 LDAP server for storage of all client information. 389 has solid replication capabilities,
but it's data sets depend on IP addresses that stay coupled to specific instance data. This creates challenges in an ephemeral 
container environment. 

The FreeIPA operator works to bridge these challenges so users can focus on the value proposition of FreeIPA instead
of spending weeks getting it settled and reliable in a continerized environment. 

## Tooling and Dependencies

The operator is based on Kubebuilder and Bazel. The container itself comes from the FreeIPA Container project. 

## Status

This project is in development phase. We'll keep this page updated as we get usable features.