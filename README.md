### Godrop
Discover peers in the local network and share files with them.
Godrop utilizes mDNS to discover the local address of potential peers. 

#### Usage
Run `godrop init` to initialize godrop and create a `.godrop.yaml` file in your home directory. You can press enter to accept all the defaults. 

### Share a file via mDNS
Run `godrop share mdns <path to file>` to share a file. On a different machine in the same network run `godrop accept ` to receive the file shared by the other machine. 