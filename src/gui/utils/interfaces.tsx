export interface IFile {
  StorageDirectoryID: number;
  Hash: string;
  ID: number;
  Name: string;
  Size: number;
}

export interface IDirectory {
  ID: number;
  Name: string;
  FileCount: number;
}

export interface ICurrentDirectory {
  ID: number;
  Name: string;
  Path: string;
  ParentDirectoryID?: number;
  PathHash: string;
}
