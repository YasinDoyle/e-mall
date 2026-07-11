import request from "@/utils/request";

export const getReviewList = (params: {
  product_id: number;
  page_num: number;
  page_size: number;
}) => request.get("/product/reviews", { params });

export const createReview = (data: {
  product_id: number;
  order_id: number;
  rating: number;
  content: string;
  images?: string;
}) => request.post("/reviews/create", data);

export const uploadReviewImage = (formData: FormData) =>
  request.post("/reviews/upload", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
