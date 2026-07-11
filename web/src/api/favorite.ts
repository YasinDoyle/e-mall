import request from "@/utils/request";

export const getFavoriteList = (params: {
  page_num: number;
  page_size: number;
}) =>
  request.get("/favorites/list", {
    params: {
      pageNum: params.page_num,
      pageSize: params.page_size,
    },
  });

export const createFavorite = (data: { product_id: number; boss_id: number }) =>
  request.post("/favorites/create", data);

export const deleteFavorite = (data: { id: number }) =>
  request.post("/favorites/delete", data);
