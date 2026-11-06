# DeviceBridge

Vocabulary for device maintenance and access. These terms distinguish the
concepts already separated in the UX and security planning documents; they
do not approve a protocol, credential mechanism or permission model.

## Language

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
