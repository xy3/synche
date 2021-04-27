import { ReactChild } from "react";
import ReactModal from "react-modal";

interface ComponentProps {
  isOpen: boolean;
  children: ReactChild | ReactChild[];
}

export default function Modal({ isOpen, children }: ComponentProps) {
  return (
    <ReactModal
      className="bg-gray-900 w-full h-full flex justify-center items-center bg-opacity-85 text-gray-400"
      isOpen={isOpen}
    >
      {children}
    </ReactModal>
  );
}
