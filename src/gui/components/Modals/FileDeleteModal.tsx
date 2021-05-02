import Modal from "../Modal";
import axios from "../../utils/axios.instance";
import cogoToast from "../../utils/cogoToast.instance";
import { RiFileWarningLine } from "react-icons/ri";

interface ComponentProps {
  fileId: number;
  isOpen: boolean;
  onSubmit(): void;
  onGoBack(): void;
}

export default function FileDeleteModal({
  fileId,
  isOpen,
  onSubmit,
  onGoBack,
}: ComponentProps) {
  async function deleteFile(fileId: number) {
    try {
      const res = await axios.delete(`/files/${fileId}`);

      if (res.status === 200) {
        onSubmit();
      }
    } catch (err) {
      cogoToast.error("There has been an error, please try again");
    }
  }

  return (
    <Modal isOpen={isOpen}>
      <div className="container">
        <h1 className="title text-gray-300">
          Are you sure you want to delete this file?
        </h1>
        <div className="my-8 p-4 bg-gray-800 shadow-sm">
          <div className="my-8 w-full flex flex-col md:flex-row">
            <div className="w-full md:w-1/2 p-2">
              <button
                className="primary-button flex justify-center items-center"
                type="button"
                onClick={() => deleteFile(fileId)}
              >
                <RiFileWarningLine className="icon" />
                <span>Yes, do it</span>
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
        </div>
      </div>
    </Modal>
  );
}
