export const isValidImageFile = (file: File | null): boolean =>
  file !== null &&
  (file.type.match("image/jpeg") !== null ||
    file.type.match("image/png") !== null);

export interface CompressionMethodConfig {
  label: string;
  maxComponents: number;
  defaultComponents: number;
  description: string;
}

export interface CompressionMethods {
  DCT: CompressionMethodConfig;
  SVD: CompressionMethodConfig;
}

export const COMPRESSION_METHODS: CompressionMethods = {
  DCT: {
    label: "DCT (Discrete Cosine Transform)",
    maxComponents: 20,
    defaultComponents: 10,
    description:
      "DCT: Block-based compression, very efficient per component (1-20)",
  },
  SVD: {
    label: "SVD (Singular Value Decomposition)",
    maxComponents: 256,
    defaultComponents: 64,
    description: "SVD: Matrix decomposition, requires more components (1-256)",
  },
};
