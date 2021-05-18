import Modal from "../Modal";
import {ErrorMessage, Field, Form, Formik} from "formik";
import axios from "../../utils/axios.instance";
import cogoToast from "../../utils/cogoToast.instance";
import * as yup from "yup";
import {RiFolderAddLine} from "react-icons/ri";

interface ComponentProps {
  currentPathID?: number;
  isOpen: boolean;
  onSubmit(): void;
  onGoBack(): void;
}

export default function FileUploadModal({
  currentPathID,
  isOpen,
  onSubmit,
  onGoBack,
}: ComponentProps) {
  const validationSchema = yup.object().shape({
    filePath: yup.string().required("File path is required"),
  });

  return (
    <Modal isOpen={isOpen}>
      <div className="container">
        <h1 className="title text-gray-300">File Upload</h1>
        <div className="my-8 p-4 bg-gray-800 shadow-sm">
          <Formik
            initialValues={{ filePath: "" }}
            onSubmit={async (values) => {
              try {
                var newFileParams = {
                  directoryID: 0,
                  filePath: values.filePath
                }

                if (currentPathID) {
                  newFileParams.directoryID = currentPathID;
                }

                const res = await axios({url: "/upload", baseURL: process.env.NEXT_PUBLIC_CLIENT_BASE_URL, method: "post", data: newFileParams});

                cogoToast.info("Uploading the file via the command line utility...")

                if (res.status === 200) {
                  onSubmit();
                }
              } catch (err) {
                cogoToast.error("There has been an error, please try again");
              }
            }}
            validationSchema={validationSchema}
          >
            {(props) => (
              <Form>
                <div className="my-4">
                  <label className="block text-gray-400 text-sm">File *</label>
                  <Field
                      type="text"
                      name="filePath"
                      className="input-inverted"
                      placeholder="Provide the full path to a file..."
                  />
                  <div className="mb-2 text-red-500 text-sm">
                    <ErrorMessage name="file" />
                  </div>
                </div>

                <div className="my-8 w-full flex flex-col md:flex-row">
                  <div className="w-full md:w-1/2 p-2">
                    <button
                      className="primary-button flex justify-center items-center"
                      type="submit"
                    >
                      <RiFolderAddLine className="icon" />
                      <span>Create</span>
                    </button>
                  </div>
                  <div className="w-full md:w-1/2 p-2">
                    <button
                      className="danger-button flex justify-center items-center"
                      type="button"
                      onClick={() => onGoBack()}
                    >
                      <span>Cancel</span>
                    </button>
                  </div>
                </div>
              </Form>
            )}
          </Formik>
        </div>
      </div>
    </Modal>
  );
}
