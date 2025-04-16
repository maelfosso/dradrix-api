import { useCallback, useMemo } from "react";
import { useLocation, useSearchParams } from "react-router-dom";

export const usePathUtils = () => {
  const location = useLocation();
  const [searchParams, setSearchParams] = useSearchParams();

  const findParameter = useCallback(
    (parameterName: string) => searchParams.get(parameterName),
    [searchParams],
  );

  const removeParameter = useCallback(
    (parameterName: string) => {
      searchParams.delete(parameterName);
      return setSearchParams(searchParams);
    },
    [searchParams, setSearchParams],
  );

  const removeParameters = useCallback(
    (parameterNames: string[]) => {
      parameterNames.forEach((parameterName) => {
        searchParams.delete(parameterName);
      });
      return setSearchParams(searchParams);
    },
    [searchParams, setSearchParams],
  );

  const addParameters = useCallback(
    (paramHash: Record<string, string>) => {
      Object.keys(paramHash).forEach((key) => {
        searchParams.set(key, paramHash[key]);
      });
      return setSearchParams(searchParams);
    },
    [searchParams, setSearchParams],
  );

  const extractParams = useCallback((path: string) => {
    const [, existingParams] = path.split("?");
    return existingParams
      ? existingParams
          .split("&")
          .reduce((hash: Record<string, string>, pair) => {
            const [key, value] = pair.split("=");
            hash[key] = decodeURIComponent(value);
            return hash;
          }, {})
      : {};
  }, []);

  const appendCurrentParams = useCallback(
    (path: string) => {
      const currentParams = extractParams(window.location.href);
      const pathParams = extractParams(path);
      const [basePath] = path.split("?");

      // eslint-disable-next-line guard-for-in
      for (const key in currentParams) {
        pathParams[key] = currentParams[key];
      }

      return `${basePath}?${Object.entries(pathParams)
        .map(([key, value]) => {
          return `${key}=${encodeURIComponent(value)}`;
        })
        .join("&")}`;
    },
    [location, extractParams],
  );

  return useMemo(
    () => ({
      findParameter,
      removeParameter,
      addParameters,
      appendCurrentParams,
      removeParameters,
    }),
    [
      appendCurrentParams,
      findParameter,
      removeParameter,
      addParameters,
      removeParameters,
    ],
  );
};