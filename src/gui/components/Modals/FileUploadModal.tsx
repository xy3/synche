import Modal from "../Modal";
import { Formik, Form, Field, ErrorMessage } from "formik";
import axios from "../../utils/axios.instance";
import cogoToast from "../../utils/cogoToast.instance";
import * as yup from "yup";
import { RiFolderAddLine } from "react-icons/ri";

interface ComponentProps {
  isOpen: boolean;
  onSubmit(): void;
  onGoBack(): void;
}

export default function FileUploadModal({
  isOpen,
  onSubmit,
  onGoBack,
}: ComponentProps) {
  const validationSchema = yup.object().shape({
    file: yup
      .mixed()
      .test("fileSize", "File Size is too large (max: 5 MB)", (value) => {
        return value ? value.size <= 5242880 : false;
      })
      .required("File is required"),
    folderId: yup.string().optional(),
  });

  return (
    <Modal isOpen={isOpen}>
      <div className="container">
        <h1 className="title text-gray-300">File Upload</h1>
        <div className="my-8 p-4 bg-gray-800 shadow-sm">
          <Formik
            initialValues={{ file: null, folderId: "" }}
            onSubmit={async (values) => {
              try {
                onSubmit();
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
                  <input
                    name="photo"
                    id="photo"
                    type="file"
                    accept="*"
                    className="input-inverted"
                    placeholder="Please upload your file"
                    onChange={(e) =>
                      props.setFieldValue(
                        "file",
                        e.currentTarget.files[0],
                        true
                      )
                    }
                  />
                  <div className="mb-2 text-red-500 text-sm">
                    <ErrorMessage name="file" />
                  </div>
                </div>

                <div className="my-4">
                  <label className="block text-gray-400 text-sm">Folder</label>
                  <Field
                    as="select"
                    className="select-inverted"
                    name="folderId"
                  >
                    <option value="">No Folder</option>
                    <option value="test">Test Folder</option>
                  </Field>

                  <div className="mb-2 text-red-500 text-sm">
                    <ErrorMessage name="folderId" />
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
