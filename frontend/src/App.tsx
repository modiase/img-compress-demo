import { useState } from "react";
import ImageUploader from "./components/ImageUploader";
import CompressionVisualizer from "./components/CompressionVisualizer";
import {
  loadCompressionData,
  saveCompressionData,
  updateSelectedComponent,
  clearCompressionData,
  CompressionData,
} from "./utils/storage";

function App() {
  const [compressionData, setCompressionData] =
    useState<CompressionData | null>(() => loadCompressionData()?.data ?? null);
  const [selectedComponent, setSelectedComponent] = useState<number | null>(
    () => loadCompressionData()?.selectedComponent ?? null,
  );
  const [loading, setLoading] = useState<boolean>(false);

  const handleCompressionComplete = (data: CompressionData): void => {
    setCompressionData(data);
    setSelectedComponent(data.componentLevels.length - 1);
    saveCompressionData(data, data.componentLevels.length - 1);
  };

  const handleComponentSelect = (index: number): void => {
    setSelectedComponent(index);
    updateSelectedComponent(index);
  };

  const handleClearData = (): void => {
    setCompressionData(null);
    setSelectedComponent(null);
    clearCompressionData();
  };

  return (
    <div className="min-h-screen bg-gray-50 p-5">
      <div className="max-w-7xl mx-auto space-y-8">
        <header className="bg-white rounded-xl shadow-sm border border-gray-200 p-8 text-center">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">
            Image Compression Tool
          </h1>
          <p className="text-lg text-gray-600">
            Compress images using DCT or SVD and visualize component data
          </p>
        </header>

        <div className="space-y-8">
          <ImageUploader
            onCompressionComplete={handleCompressionComplete}
            onCompressionStart={handleClearData}
            setLoading={setLoading}
          />

          {loading && (
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-10 text-center">
              <div className="w-12 h-12 border-4 border-gray-200 border-t-primary rounded-full animate-spin mx-auto mb-5"></div>
              <p className="text-gray-700">Compressing image...</p>
            </div>
          )}

          {compressionData && !loading && (
            <CompressionVisualizer
              compressionData={compressionData}
              selectedComponent={selectedComponent}
              onComponentSelect={handleComponentSelect}
              onClearResults={handleClearData}
            />
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
