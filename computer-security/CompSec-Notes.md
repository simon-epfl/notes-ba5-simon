OS --> interface entre l'utilisateur et le hardware

Complete mediation --> all requests are mediated—all requests go to
the reference monitor that enforces specified access control policies

The reference monitor grants permission to users to apply certain
operations to a given resource

Execute permission on a directory allows traversing it
Read permission on a directory allows lookup

► When a user runs a process, it runs with that user’s privileges, i.e.

they can access any resource that user has permissions for

► By default, a child process inherits its parent’s privileges


Every process has:

► Real user ID (uid) - the user ID that started that process

► Effective user ID (euid) - the user ID that determines the

process’ privileges

► Saved user ID (suid) - the effective user ID before the last

Normally, when you execute a program:

- The process runs with **your** user ID (RUID = EUID = your UID).

But if the program’s **setuid bit** is set, then:

- The process’s **effective user ID (EUID)** becomes the _file owner’s UID_.

That’s what allows certain trusted programs to perform sensitive operations safely.




modification

- Most processes have EUID = RUID, so they have the same privileges as the user who ran them.
- But in special cases — for example, **setuid programs** — the EUID can differ.
![[image-32.png]]

but permissions are too coarse-grained:

Because you can’t easily express more nuanced or specific access rules.

**Examples of situations the model can’t handle well:**

1. _“Alice and Bob can read the file, but Charlie and Dana can’t; however, Eve can write to it.”_  
    → To do that, you’d have to make a _new group_ with exactly Alice and Bob in it, then assign it to the file. Clunky.
    
2. _“Finance can read; IT can read and write.”_  
    → Only one group can be attached to a file, so you can’t assign per‑group permissions directly.
    
3. _“Give user X read access for 24 hours, then revoke automatically.”_  
    → Impossible with standard rwx flags.
    

These all require more precise (fine‑grained) control than the model offers.


![[image-33.png]]
![[image-34.png]]