import { isSuccessResponse, sendRequest } from '../../httpClient';
import { UserContext } from '../../types';

export const setEntities = async <T>(
  ctx: UserContext,
  entities: T[],
  endpoint: string,
  setter: (body: { [key: string]: T[] }) => void
) => {
  const res = await sendRequest(ctx, 'GET', endpoint);
  if (!isSuccessResponse(res)) {
    entities.length = 0;
    return;
  }
  const body: { [key: string]: T[] } = await res.json();
  setter(body);
};

export const withPromptValues = (fn: (...values: string[]) => void, ...prompts: string[]) => {
  const args = [];
  for (const p of prompts) {
    const answer = prompt(p);
    if (answer === null) return;
    args.push(answer);
  }
  fn(...args);
};
