export interface ComponentLevel {
  numComponents: number;
  dataSize: number;
  imageData: string;
}

export interface CompressionData {
  method: string;
  originalSize: number;
  componentLevels: ComponentLevel[];
}

export interface SavedCompressionData {
  data: CompressionData;
  selectedComponent: number | null;
}

const STORAGE_KEY = "img-compress-data";
const SELECTED_KEY = "img-compress-selected";

export const saveCompressionData = (
  data: CompressionData,
  selectedIndex: number,
): void => {
  try {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(data));
    sessionStorage.setItem(SELECTED_KEY, selectedIndex.toString());
  } catch (error) {
    console.error("Failed to save to session storage:", error);
  }
};

export const loadCompressionData = (): SavedCompressionData | null => {
  try {
    const savedData = sessionStorage.getItem(STORAGE_KEY);
    if (!savedData) return null;
    const data = JSON.parse(savedData) as CompressionData;
    const savedSelected = sessionStorage.getItem(SELECTED_KEY);
    return {
      data,
      selectedComponent:
        savedSelected !== null
          ? parseInt(savedSelected)
          : data.componentLevels.length > 0
            ? data.componentLevels.length - 1
            : null,
    };
  } catch (error) {
    console.error("Failed to load from session storage:", error);
    clearCompressionData();
    return null;
  }
};

export const updateSelectedComponent = (index: number): void => {
  try {
    sessionStorage.setItem(SELECTED_KEY, index.toString());
  } catch (error) {
    console.error("Failed to update session storage:", error);
  }
};

export const clearCompressionData = (): void => {
  sessionStorage.removeItem(STORAGE_KEY);
  sessionStorage.removeItem(SELECTED_KEY);
};
