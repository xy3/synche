# Proposal

# School of Computing - Year 4 Project Proposal Form

## SECTION A

|Project Title:       | Synche            |
|---------------------|-------------------|
|Student 1 Name:      | Tara Collins      |
|Student 1 ID:        | 15416118          |
|Student 2 Name:      | Theo Coyne Morgan |
|Student 2 ID:        | 17338811          |
|Project Supervisor:  | Dr. Stephen Blott |

<img src='https://i.postimg.cc/QCRDPv9y/new-crop.png' width='400'/>

## SECTION B

### Introduction

The project idea is a multipart file uploading tool and storage system. It will be developed as a command line tool (and/or desktop GUI application) that splits files into small parts and concurrently uploads the parts to a storage server. The software for the storage server is included in the project.

It is possible to expand this idea further and develop a full web application that can automatically deploy the required software to a user's server.

### Outline

The proposed project idea is a multipart upload tool which will segment files and use concurrency to upload these segmented parts to a private server where they are coalesced into a composite file.

There are multiple parts to this project. There is the client side that will split, upload, compress, and encrypt files. It will also provide authentication and manage the file segmenter.

The server side will reassemble the file and manage the concurrent uploading of the file segments. It will also manage the file reassembler.

Lastly, it will use a database for file tracking and file hash storage.

![Usage diagram](https://i.imgur.com/f5MAWO7.png)

### Background

The original idea was to create an all-in-one system which would integrate with popular cloud storage services e.g. Google Drive, and support multipart file uploading to increase upload speed and reliability. Uploading large files or multiple files using these services can be slow and a broken connection may result in a full failure to upload. Upon further research, we discovered that creating a universal tool for these existing services would not be possible because there is limited access to their APIs.

We were still interested in creating some form of upload tool when we discovered Amazon S3's (Amazon Web Services) multipart uploading tool. As the name suggests, with Amazon S3, the user may send a request to initiate a multipart upload which uploads multiple parts of the same file concurrently instead of sequentially. 

Amazon S3 works by returning a response with an upload ID after the user has requested to initiate the multipart upload. This unique ID is used to upload parts, list the parts, complete an upload, or stop an upload.  Amazon S3 expects that the user splits their files themselves and numbers each part so the multipart uploading tool knows the order that the file should be reassembled. 

We thought that the AWS tool is incredibly useful but it is not accessible. AWS products are expensive and cater for large businesses rather than an individual user. The AWS tool also does not provide a mechanism for splitting large files, it expects that the user does that themselves.

We decided that we wanted to create a tool that is similar to the Amazon S3 multipart tool but is more accessible and user friendly. Currently, a full package multipart upload tool does not exist. This is what we aim to create.

### Achievements

The aim of Synche is to provide a faster and more reliable alternative to cloud storage solutions such as Google Drive. It will achieve this by making use of multipart data transfer and (self-served) server software for storing files.

Current commercial cloud storage does not offer the general public a method to upload files concurrently in multiple parts.

The users will range from individuals that may use the tool for personal use to small-medium companies for internal use.

**Universal benefits of using Synche are:**

- **Faster file transfer**  
By using multipart uploading, higher data transfer throughput can be achieved for large files (>100MB).
- **Better error recovery**  
    Files are uploaded in small parts, so if the connection fails, the parts that have been uploaded already are still there.
- **Pausing and resuming uploads**  
    Uploads can be paused and the non-transferred chunks can be uploaded at another time.
- **Full ownership of the server**  
    Cloud storage such as OneDrive etc. may be a privacy or security concern for some people, and by owning your own server you can be in full control.
- **No proprietary software e.g. Google Drive**  
    The software is open-source, meaning there are no fees, and it is not tied to any company.
- **Cheaper storage costs**  
    Storage solutions such as Amazon S3 can become extremely costly for huge amounts of data storage. By owning your own server you can dramatically cut down the storage costs.

### Justification

Outlined in this section, are two distinct use cases where users would benefit from a multipart file transfer and self-hosted storage solution.

**For a digital agency**

A small to medium company such as a photography or film agency would be a perfect use case for our project. The company might need to transfer files quickly to a central storage location, such as the company's own servers. For instance, a film crew may want to upload their footage to a shared server. 

*Example use case:*

An international film crew captures 500GB of footage between them. They want to store the footage on a shared server so the editing team can start working on it. Each member is in a different country, with varying connection speeds. 

By using our software, the team can get the footage to the editors faster, avoid having to re-upload the files if the connection drops, and stop and start the transfer at their leisure. The proposed software is particularly useful if the connection where the files are being uploaded from is not guaranteed to be fast or stable.

**As a cloud alternative**

An independent user has sensitive or private information that they want to upload to their own private server. Their main concern is the security of their information.

*Example use case:*

A user has private documents that they want to back-up on a server that they already own. They do not want to use a service such a Google Drive as they currently pay for their own server and want to ensure that their information is truly private.

By using the proposed software, the user has full ownership over their data. If the user lives in a location with a medium-slow internet connection, they may transfer their data quicker and not need to be as concerned about loss of data in the event of a connection failure.

### Programming language(s)

*  Golang
*  Javascript
*  PHP
*  Clojure*

\* Not guaranteed to be used

### Programming tools / Tech stack

- Server API
    - Golang HTTP Routing
    - Nginx Web server

- Client / Desktop app
    - go-astilectron (Electron framework for Golang)
    - HTML / SCSS / JS (Frontend design)

- File Database
    - MongoDB or Cassandra
    - Redis for caching

- Web dashboard
    - Laravel
    - HTML / SCSS / JS (Frontend design)
    - Grafana (Data visualization)

- Testing
    - Table driven tests (Advanced Golang testing)
    - Httpexpect (API testing Golang library)

- Project management
    - Clubhouse.io project tracker
    - Gitlab Issues bug tracker

- Development
    - GoLand IDE or Visual Studio Code

### Hardware

- A Linux storage server for testing

### Learning Challenges

> technologies, languages, tools

- Golang - both project members are relatively new to Go
- MongoDB and Redis caching
- Security and privacy best practices for data storage
- API testing
- API security
- Storage allocation

Additional potential learning challenges
- Grafana integration

### Breakdown of work

> It must be clear from the explanation of this breakdown of work both that each student is responsible for
> separate, clearly-defined tasks, and that those responsibilities substantially cover all of the work required
> for the project.

We intend to both work on all aspects of the project, as pair programming often results in better quality code.

We both especially intend to work on parallel uploading as this is the core feature of our project.

#### Student 1

Tara:

- Database implementation
- Client command line
- File serving


#### Student 2

Theo:

- File segmentation
- Data management
- Server API
