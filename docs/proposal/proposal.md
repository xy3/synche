# Proposal

# School of Computing &mdash; Year 4 Project Proposal Form

> Edit (then commit and push) this document to complete your proposal form.
> Make use of figures / diagrams where appropriate.
>
> Do not rename this file.

## SECTION A

|                     |                   |
|---------------------|-------------------|
|Project Title:       | Synche            |
|Student 1 Name:      | Tara Collins      |
|Student 1 ID:        | 15416118          |
|Student 2 Name:      | Theo Coyne Morgan |
|Student 2 ID:        | 17338811          |
|Project Supervisor:  | Dr. Stephen Blott |

> Ensure that the Supervisor formally agrees to supervise your project; this is only recognised once the
> Supervisor assigns herself/himself via the project Dashboard.
>
> Project proposals without an assigned
> Supervisor will not be accepted for presentation to the Approval Panel.

## SECTION B

> Guidance: This document is expected to be approximately 3 pages in length, but it can exceed this page limit.
> It is also permissible to carry forward content from this proposal to your later documents (e.g. functional
> specification) as appropriate.
>
> Your proposal must include *at least* the following sections.


### Introduction

> Describe the general area covered by the project.

Synche will be a composite uploading tool for managing file uploads to a storage server. It will be developed as a command line tool that splits files and uses concurrency to upload the file parts to the storage server. It is also possible to expand this tool further and make it a full web application that also deploys the user's server.

### Outline

> Outline the proposed project.

![Usage diagram](https://i.imgur.com/yGoBU4a.png)

Synche be a multipart upload tool which will segment files and use concurrency to upload these segmented parts to a private server where they are coalesced into a composite file.

There are multiple parts to this project. There is the client side that will split, upload, compress, and encrypt files. It will also provide authentication and manage the file segmenter.

The server side will reassemble the file and manage the concurrent uploading of the file segments. It will also manage the file reassembler.

Lastly, it will use a database for file tracking and file hash storage.

### Background

> Where did the ideas come from?

The original idea was to create an all-in-one system which would integrate with popular cloud storage services e.g. Google Drive, and support multipart file uploading to increase upload speed and reliability. Uploading large files or multiple files using these services can be slow and a broken connection may result in a full failure to upload. Upon further research, we discovered that creating a universal tool for these existing services would not be possible because there is limited access to their APIs.

We were still interested in creating some form of upload tool when we discovered Amazon S3's (Amazon Web Services) multipart uploading tool. As the name suggests, with Amazon S3, the user may send a request to initiate a multipart upload which uploads multiple parts of the same file concurrently instead of sequentially. 

Amazon S3 works by returning a response with an upload ID after the user has requested to initiate the multipart upload. This unique ID is used to upload parts, list the parts, complete an upload, or stop an upload.  Amazon S3 expects that the user splits their files themselves and numbers each part so the multipart uploading tool knows the order that the file should be reassembled. 

We thought that the AWS tool is incredibly useful but it is not accessible. AWS products are expensive and cater for large businesses rather than an individual user. The AWS tool also does not provide a mechanism for splitting large files, it expects that the user does that themselves.

We decided that we wanted to create a tool that is similar to the Amazon S3 multipart tool but is more accessible and user friendly. Currently, a full package multipart upload tool does not exist. This is what we aim to create.

### Achievements

> What functions will the project provide? Who will the users be?

The users will range from individuals that may use the tool for personal use and small companies that may use the tool internally.

### Justification

> Why/when/where/how will it be useful?


There are many scenarios in which Synche will be useful. The universal benefits of using our tool will be:

- **Faster file transfer**
    By using multipart uploading, higher transfer rates can be achieved for large files (>100MB)
- **Better error recovery**
    Files are uploaded in small parts, so if the connection fails, the parts that have been uploaded already are still there
- **Pausing and resuming uploads**
    Uploads can be paused and the non-transferred chunks can be uploaded at another time
- **Full ownership of the server**
    Cloud storage such as OneDrive etc. may be a privacy or security concern for some people, and by owning your own server you can be in full control
- **No proprietary software e.g. Google Drive**
    The software is open-source, meaning there are no fees, and it is not tied to any company
- **Cheaper storage costs**
    Storage solutions such as Amazon S3 can become extremely costly for huge amounts of data storage. By owning your own server you can dramatically cut down the storage costs.

Below we have outlines various scenarios where the users could greatly benefit from multipart file uploading and a self-served storage solution. These examples are:

#### **The digital agency**

A small to medium company such as a photography or film agency would be a perfect use case for our project. The company might need to transfer files quickly to a central storage location, such as the company's own servers. For instance, a film crew may want to upload their footage to a shared server. 

#### *Example use case:*

An international film crew captures 500GB of footage between them. They want to store the footage on a shared server so the editing team can start working on it. Each member is in a different country, with varying connection speeds. 

By using our software, the team can get the footage to the editors faster, avoid having to re-upload the files if the connection drops, and stop and start the transfer at their leisure. The proposed software is particularly useful if the connection where the files are being uploaded from is not guaranteed to be fast or stable.

#### **The video server owner**

#### **The cloud alternative**
An independent user has sensitive or private information that they want to upload to their own private server. Their main concern is the security of their information.

#### *Example use case:*
A user has private documents that they want to back up on a server that they already own. They do not want to use a service such a Google Drive as they already pay for their own server and they want to ensure that their information is secure.

By using our software, the user has full ownership over their data. If they live somewhere with slow internet speeds, they may transfer their data quicker and not need to be concerned about loss of data.

### Programming language(s)
> List the proposed language(s) to be used.

*  Golang
*  Java
*  Clojure 

### Programming tools / Tech stack

> Describe the compiler, database, web server, etc., and any other software tools you plan to use.

- File database: (MongoDB or Cassandra) + Redis for caching
- Server API - Golang HTTP Routing
- Desktop app - go-astilectron (Electron)
- HTML/CSS/JS for the UI
- Testing - Table driven tests (for Golang)
- Concurrency - Clojure

### Hardware

> Describe any non-standard hardware components which will be required.

There will be no non-standard hardware components required to complete this project.

### Learning Challenges

> List the main new things (technologies, languages, tools, etc) that you will have to learn.

We both will be learning Go as we work on this project. Currently, we have basic knowledge of the language and we have only written basic scripts in it. We plan to use Go for the majority of this application. 

### Breakdown of work

> Clearly identify who will undertake which parts of the project.
>
> It must be clear from the explanation of this breakdown of work both that each student is responsible for
> separate, clearly-defined tasks, and that those responsibilities substantially cover all of the work required
> for the project.


We intend to both work on all aspects of the project, as pair programming often results in better quality code.

We both especially intend to work on parallel uploading as this is the core feature of our project.

#### Student 1

> *Student 1 should complete this section.*

Tara:

- Database implementation
- Client command line
- File serving


#### Student 2

> *Student 2 should complete this section.*

Theo:

- File segmentation
- Data management
- Server API

## Example

> Example: Here's how you can include images in markdown documents...
