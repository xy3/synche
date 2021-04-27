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

export default function NewFolderModal({
  isOpen,
  onSubmit,
  onGoBack,
}: ComponentProps) {
  const validationSchema = yup.object().shape({
    folderName: yup.string().required("Folder Name is required"),
  });

  return (
    <Modal isOpen={isOpen}>
      <div className="container">
        <h1 className="title text-gray-300">New Folder</h1>
        <div className="my-8 p-4 bg-gray-800 shadow-sm">
          <Formik
            initialValues={{ folderName: "" }}
            onSubmit={async (values) => {
              try {
                onSubmit();
              } catch (err) {
                cogoToast.error("There has been an error, please try again");
              }
            }}
            validationSchema={validationSchema}
          >
            <Form>
              <div className="my-4">
                <label className="block text-gray-400 text-sm">
                  Folder Name *
                </label>
                <Field
                  type="text"
                  name="folderName"
                  className="input-inverted"
                  placeholder="Please enter new's folder name..."
                />
                <div className="mb-2 text-red-500 text-sm">
                  <ErrorMessage name="folderName" />
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
          </Formik>
        </div>
      </div>
    </Modal>
  );
}
