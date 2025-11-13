export const formatBytes = (bytes: number): string => (bytes / 1024).toFixed(2);

export const formatCompressionRatio = (
  originalSize: number,
  compressedSize: number,
): string => (originalSize / compressedSize).toFixed(2);

export const formatSizePercentage = (
  size: number,
  originalSize: number,
): string => ((size / originalSize) * 100).toFixed(1);

export const formatSizePerComponent = (
  dataSize: number,
  numComponents: number,
): string => (dataSize / numComponents / 1024).toFixed(2);

export const pluralize = (
  count: number,
  singular: string = "",
  plural: string = "s",
): string => (count === 1 ? singular : plural);
