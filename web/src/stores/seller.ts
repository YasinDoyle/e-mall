import { computed, ref } from "vue";
import { defineStore } from "pinia";
import { getSellerProfile, type SellerProfile } from "@/api/seller";

const SELLER_PROFILE_CACHE_KEY = "sellerProfile";

export const useSellerStore = defineStore("seller", () => {
  const profile = ref<SellerProfile | null>(
    JSON.parse(localStorage.getItem(SELLER_PROFILE_CACHE_KEY) ?? "null"),
  );
  const loaded = ref(!!profile.value);
  const loading = ref(false);

  const isApproved = computed(() => profile.value?.status === 1);
  const hasApplied = computed(() => !!profile.value);

  function setProfile(nextProfile: SellerProfile | null) {
    profile.value = nextProfile;
    loaded.value = true;
    if (nextProfile) {
      localStorage.setItem(
        SELLER_PROFILE_CACHE_KEY,
        JSON.stringify(nextProfile),
      );
    } else {
      localStorage.removeItem(SELLER_PROFILE_CACHE_KEY);
    }
  }

  async function loadProfile(options?: { force?: boolean; silentError?: boolean }) {
    if (loaded.value && !options?.force) {
      return profile.value;
    }
    loading.value = true;
    try {
      const res: any = await getSellerProfile({
        silentError: options?.silentError ?? true,
      });
      setProfile(res.data ?? null);
      return profile.value;
    } catch {
      setProfile(null);
      return null;
    } finally {
      loading.value = false;
    }
  }

  function clearProfile() {
    setProfile(null);
    loaded.value = false;
  }

  return {
    profile,
    loaded,
    loading,
    isApproved,
    hasApplied,
    setProfile,
    loadProfile,
    clearProfile,
  };
});
