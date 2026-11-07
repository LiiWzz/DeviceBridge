# DeviceBridge

Vocabulary for device maintenance and access. These terms distinguish the
concepts already separated in the UX and security planning documents; they
do not approve a protocol, credential mechanism or permission model.

## Language

**Operator identity**:
The distinct human principal responsible for a maintenance action within an
authoritative identity domain. Its recognition does not itself grant permission.
_Avoid_: Display name, device identity, SSH account

**Maintenance authority**:
Permission for an operator to perform specified maintenance actions on specified
resources within an authoritative scope.
_Avoid_: Online status, authenticated identity, unrestricted administrator

**Device record**:
The administrative entry describing an intended device. Its presence or display
metadata does not establish which device is present or whether it can be reached.
_Avoid_: Enrolled device, verified identity

**Device identity**:
The distinct security principal representing an agent's device. Recognition of
that identity is separate from its current credential validity or access rights.
_Avoid_: Device name, inventory record

**Enrollment**:
The authorized establishment of a binding between an intended device record and
a verified device identity.
_Avoid_: Record creation, installation, online status

**Intended-device proof**:
Evidence connecting the enrollment claimant to the independently established
identity of the device intended for a particular record.
_Avoid_: Record identifier, enrollment permission, physical authenticity

**Connectivity observation**:
Evidence about a device's reachability at an observed time. It does not establish
enrollment or permission to open an SSH session.
_Avoid_: Enrollment success, SSH authorization

**Session authorization**:
Permission for a particular user to request access to a particular device in a
particular session context.
_Avoid_: Inventory visibility, enrollment authority

**Browser terminal**:
An interactive command-line view inside the dashboard for an authorized device
session. The view itself grants no device or account permissions.
_Avoid_: One-command form, local computer shell
