import {
  formatBytes,
  formatCompressionRatio,
  formatSizePercentage,
  formatSizePerComponent,
  pluralize,
} from "../utils/formatting";
import { CompressionData } from "../utils/storage";

interface CompressionVisualizerProps {
  compressionData: CompressionData | null;
  selectedComponent: number | null;
  onComponentSelect: (index: number) => void;
  onClearResults: () => void;
}

function CompressionVisualizer({
  compressionData,
  selectedComponent,
  onComponentSelect,
  onClearResults,
}: CompressionVisualizerProps) {
  if (!compressionData || selectedComponent === null) return null;
  const currentLevel = compressionData.componentLevels[selectedComponent];

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <div className="flex justify-between items-center mb-6 flex-wrap gap-4">
        <h2 className="text-2xl font-semibold text-gray-800">
          Compression Results
        </h2>
        <div className="flex items-center gap-3">
          <div className="bg-primary px-5 py-2 rounded-full font-semibold text-gray-800 text-sm">
            {compressionData.method}
          </div>
          <button
            onClick={onClearResults}
            className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium text-sm transition-colors"
            title="Clear results and start over"
          >
            Clear Results
          </button>
        </div>
      </div>

      <div className="mb-8">
        <div className="mb-5">
          <h3 className="text-xl text-gray-800 font-medium mb-4">
            Reconstructed with {currentLevel.numComponents} component
            {pluralize(currentLevel.numComponents)}
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="flex flex-col gap-1">
              <span className="text-sm text-gray-600 font-medium">
                Data Size:
              </span>
              <span className="text-2xl text-primary-dark font-bold">
                {formatBytes(currentLevel.dataSize)} KB
              </span>
              <span className="text-xs text-gray-500">
                {formatSizePercentage(
                  currentLevel.dataSize,
                  compressionData.originalSize,
                )}
                % of original
              </span>
            </div>
            <div className="flex flex-col gap-1">
              <span className="text-sm text-gray-600 font-medium">
                Compression Ratio:
              </span>
              <span className="text-2xl text-primary-dark font-bold">
                {formatCompressionRatio(
                  compressionData.originalSize,
                  currentLevel.dataSize,
                )}
                :1
              </span>
              <span className="text-xs text-gray-500">
                {formatCompressionRatio(
                  compressionData.originalSize,
                  currentLevel.dataSize,
                )}
                x smaller
              </span>
            </div>
            <div className="flex flex-col gap-1">
              <span className="text-sm text-gray-600 font-medium">
                Size per Component:
              </span>
              <span className="text-2xl text-primary-dark font-bold">
                {formatSizePerComponent(
                  currentLevel.dataSize,
                  currentLevel.numComponents,
                )}{" "}
                KB
              </span>
              <span className="text-xs text-gray-500">avg. per component</span>
            </div>
          </div>
        </div>

        <div className="bg-gray-50 rounded-lg p-5 flex justify-center items-center min-h-[400px]">
          <img
            src={`data:image/png;base64,${currentLevel.imageData}`}
            alt="Compressed"
            className="max-w-full max-h-[600px] object-contain rounded shadow-md"
          />
        </div>
      </div>

      <div className="flex flex-col gap-2">
        <label htmlFor="component-slider" className="font-medium text-gray-800">
          Select Component Level: <strong>{currentLevel.numComponents}</strong>{" "}
          / {compressionData.componentLevels.length}
        </label>
        <input
          id="component-slider"
          type="range"
          min="0"
          max={compressionData.componentLevels.length - 1}
          value={selectedComponent}
          onChange={(e) => onComponentSelect(parseInt(e.target.value))}
          className="w-full h-2 bg-gradient-to-r from-primary to-secondary rounded-lg appearance-none cursor-pointer"
          style={{ WebkitAppearance: "none" }}
        />
        <style jsx>{`
          input[type="range"]::-webkit-slider-thumb {
            -webkit-appearance: none;
            appearance: none;
            width: 24px;
            height: 24px;
            border-radius: 50%;
            background: white;
            border: 3px solid #9bc0e8;
            cursor: pointer;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
          }
          input[type="range"]::-webkit-slider-thumb:hover {
            transform: scale(1.2);
          }
          input[type="range"]::-moz-range-thumb {
            width: 24px;
            height: 24px;
            border-radius: 50%;
            background: white;
            border: 3px solid #9bc0e8;
            cursor: pointer;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
          }
          input[type="range"]::-moz-range-thumb:hover {
            transform: scale(1.2);
          }
        `}</style>
        <div className="flex justify-between text-sm text-gray-500">
          <span>1 component</span>
          <span>{compressionData.componentLevels.length} components</span>
        </div>
      </div>
    </div>
  );
}

export default CompressionVisualizer;
