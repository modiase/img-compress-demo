import { useState, useRef } from "react";
import axios from "axios";
import clsx from "clsx";
import { isValidImageFile, COMPRESSION_METHODS } from "../utils/validation";
import { CompressionData } from "../utils/storage";

interface ImageUploaderProps {
  onCompressionComplete: (data: CompressionData) => void;
  onCompressionStart?: () => void;
  setLoading: (loading: boolean) => void;
}

type CompressionMethod = "DCT" | "SVD";

function ImageUploader({
  onCompressionComplete,
  onCompressionStart,
  setLoading,
}: ImageUploaderProps) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [method, setMethod] = useState<CompressionMethod>("DCT");
  const [numComponents, setNumComponents] = useState<number>(10);
  const [preview, setPreview] = useState<string | null>(null);
  const [dragActive, setDragActive] = useState<boolean>(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleMethodChange = (newMethod: CompressionMethod): void => {
    setMethod(newMethod);
    if (
      numComponents > COMPRESSION_METHODS[newMethod].maxComponents ||
      (numComponents <= COMPRESSION_METHODS[method].maxComponents &&
        numComponents < COMPRESSION_METHODS[newMethod].maxComponents / 3)
    ) {
      setNumComponents(COMPRESSION_METHODS[newMethod].defaultComponents);
    }
  };

  const handleFileSelect = (file: File | null): void => {
    if (!isValidImageFile(file))
      return alert("Please select a JPG or PNG image");
    setSelectedFile(file);
    const reader = new FileReader();
    reader.onloadend = () => setPreview(reader.result as string);
    reader.readAsDataURL(file);
  };

  const handleDrag = (e: React.DragEvent<HTMLDivElement>): void => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(e.type === "dragenter" || e.type === "dragover");
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>): void => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    if (e.dataTransfer.files?.[0]) handleFileSelect(e.dataTransfer.files[0]);
  };

  const handleCompress = async (): Promise<void> => {
    if (!selectedFile) return alert("Please select an image first");
    onCompressionStart?.();
    setLoading(true);
    const formData = new FormData();
    formData.append("image", selectedFile);
    formData.append("method", method);
    formData.append("numComponents", numComponents.toString());
    try {
      onCompressionComplete(
        (
          await axios.post("http://localhost:8080/api/compress", formData, {
            headers: { "Content-Type": "multipart/form-data" },
          })
        ).data,
      );
    } catch (error) {
      console.error("Compression failed:", error);
      alert(
        "Compression failed: " +
          (axios.isAxiosError(error)
            ? error.response?.data?.error || error.message
            : String(error)),
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <h2 className="text-2xl font-semibold text-gray-800 mb-6">
        Upload Image
      </h2>

      <div
        className={clsx(
          "border-2 border-dashed rounded-lg p-10 text-center cursor-pointer transition-all mb-8",
          dragActive
            ? "border-primary-dark bg-blue-50 scale-[1.02]"
            : "border-gray-300 bg-gray-50",
          preview && "p-0 border-none bg-transparent",
          !dragActive && !preview && "hover:border-primary hover:bg-blue-50/50",
        )}
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
        onClick={() => fileInputRef.current?.click()}
      >
        {preview ? (
          <div className="relative rounded-lg overflow-hidden group">
            <img
              src={preview}
              alt="Preview"
              className="w-full h-[300px] object-contain bg-gray-50"
            />
            <div className="absolute inset-0 bg-black/70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
              <p className="text-white text-xl font-medium">
                Click to change image
              </p>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center gap-4 text-gray-600">
            <svg
              width="64"
              height="64"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              className="text-primary-dark"
            >
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="17 8 12 3 7 8" />
              <line x1="12" y1="3" x2="12" y2="15" />
            </svg>
            <p className="text-lg font-medium text-gray-800">
              Drop image here or click to browse
            </p>
            <span className="text-sm text-gray-400">Supports JPG and PNG</span>
          </div>
        )}
        <input
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png"
          onChange={(e) =>
            handleFileSelect(e.target.files ? e.target.files[0] : null)
          }
          className="hidden"
        />
      </div>

      <div className="space-y-6">
        <div className="flex flex-col gap-2">
          <label htmlFor="method" className="font-medium text-gray-800">
            Compression Method
          </label>
          <select
            id="method"
            value={method}
            onChange={(e) =>
              handleMethodChange(e.target.value as CompressionMethod)
            }
            className="px-3 py-2.5 border-2 border-gray-200 rounded-md bg-white text-gray-800 cursor-pointer transition-colors hover:border-primary focus:border-primary focus:outline-none"
          >
            <option value="DCT">{COMPRESSION_METHODS.DCT.label}</option>
            <option value="SVD">{COMPRESSION_METHODS.SVD.label}</option>
          </select>
          <p className="text-xs text-gray-500">
            {COMPRESSION_METHODS[method].description}
          </p>
        </div>

        <div className="flex flex-col gap-2">
          <label htmlFor="components" className="font-medium text-gray-800">
            Number of Components: <strong>{numComponents}</strong>
          </label>
          <input
            id="components"
            type="range"
            min="1"
            max={COMPRESSION_METHODS[method].maxComponents}
            value={numComponents}
            onChange={(e) => setNumComponents(parseInt(e.target.value))}
            className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-primary-dark"
          />
          <div className="flex justify-between text-sm text-gray-500">
            <span>1</span>
            <span>{COMPRESSION_METHODS[method].maxComponents}</span>
          </div>
        </div>

        <button
          onClick={handleCompress}
          disabled={!selectedFile}
          className={clsx(
            "w-full px-8 py-4 rounded-lg font-semibold text-lg transition-all",
            selectedFile
              ? "bg-primary-dark text-gray-800 hover:-translate-y-0.5 hover:shadow-lg"
              : "bg-gray-300 text-gray-500 cursor-not-allowed",
          )}
        >
          Compress Image
        </button>
      </div>
    </div>
  );
}

export default ImageUploader;
