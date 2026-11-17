import { Job, JobFilterParams, PaginatedJobsResponse } from "@/types/job";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const SAMPLE_JOBS: Job[] = [
  {
    id: "9f1e1a80-4966-419b-a01b-bf2d6e382d01",
    hash_key: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    title: "Graduate Software Engineer - Backend (Go/Python)",
    company: "Razorpay",
    role: "Backend Engineer",
    city: "Bangalore",
    source_url: "https://razorpay.com/jobs/grad-backend",
    created_at: new Date(Date.now() - 1000 * 60 * 35).toISOString(),
  },
  {
    id: "8c2d2b70-3855-408a-9f0a-ae1c5d271c02",
    hash_key: "f4c1d55309fd2d250bfca5d9007fc03538bf52f5750c045db506002c8963c966",
    title: "Frontend Intern - React & Next.js",
    company: "Swiggy",
    role: "Frontend Developer",
    city: "Bangalore",
    source_url: "https://careers.swiggy.com/jobs/frontend-intern",
    created_at: new Date(Date.now() - 1000 * 60 * 95).toISOString(),
  },
  {
    id: "7b3c3c60-2744-4f79-8e9b-9d0b4c160b03",
    hash_key: "a5b2e66410fe3e361cfdb6ea118fd14649cf63e6861d156ec617113d9074da77",
    title: "Associate DevOps Engineer (Fresher 2026 Batch)",
    company: "CRED",
    role: "DevOps",
    city: "Bangalore",
    source_url: "https://careers.cred.club/jobs/devops-fresher",
    created_at: new Date(Date.now() - 1000 * 60 * 180).toISOString(),
  },
  {
    id: "6a4d4d50-1633-4e68-7d8a-8c9a3b05fa04",
    hash_key: "b6c3f77521af4f472d0ec7fb229fe25750df74f7972e267fd728224ea185eb88",
    title: "Junior Full Stack Developer",
    company: "Postman",
    role: "Full Stack",
    city: "Bangalore",
    source_url: "https://www.postman.com/careers/jr-fullstack",
    created_at: new Date(Date.now() - 1000 * 60 * 250).toISOString(),
  },
  {
    id: "595e5e40-0522-4d57-6c79-7b892a94e905",
    hash_key: "c7d4a88632bf5a583e1fd8ac33a0f36861ea85a8a83f378ae839335fb296fc99",
    title: "Data Analyst Trainee / Fresher",
    company: "PhonePe",
    role: "Data / AI",
    city: "Pune",
    source_url: "https://www.phonepe.com/careers/data-analyst-trainee",
    created_at: new Date(Date.now() - 1000 * 60 * 360).toISOString(),
  },
  {
    id: "486f6f30-f411-4c46-5b68-6a781983d806",
    hash_key: "d8e5b99743cf6b694f2ae9bd44b1a47972fb96b9b94a489bf940446ac307ad00",
    title: "Software Engineering Intern - Systems & Core",
    company: "Zerodha",
    role: "Backend Engineer",
    city: "Bangalore",
    source_url: "https://careers.zerodha.com/interns",
    created_at: new Date(Date.now() - 1000 * 60 * 480).toISOString(),
  },
  {
    id: "377a7a20-e300-4b35-4a57-59670872c707",
    hash_key: "e9f6ca0854df7c7a5a3bf0ce55c2b58a83ac07cac05b59ac0a51557bd418be11",
    title: "Mobile App Developer Intern (Flutter / React Native)",
    company: "Groww",
    role: "Mobile",
    city: "Bangalore",
    source_url: "https://groww.in/careers/mobile-intern",
    created_at: new Date(Date.now() - 1000 * 60 * 600).toISOString(),
  },
  {
    id: "268b8b10-d2ff-4a24-3946-4856f761b608",
    hash_key: "fa07db1965e08d8b6b4c01df66d3c69b94bd18dbd16c60bd1b62668ce529cf22",
    title: "Software Engineer Trainee - Cloud Infrastructure",
    company: "Persistent Systems",
    role: "DevOps",
    city: "Pune",
    source_url: "https://careers.persistent.com/trainee",
    created_at: new Date(Date.now() - 1000 * 60 * 720).toISOString(),
  },
  {
    id: "159c9c00-c1ee-4913-2835-3745e650a509",
    hash_key: "0b18ec2a76f19e9c7c5d12e077e4d7ac05ce29ece27d71ce2c73779df63ade33",
    title: "Frontend Web Developer (0-1 Yrs Exp)",
    company: "InMobi",
    role: "Frontend Developer",
    city: "Hyderabad",
    source_url: "https://www.inmobi.com/company/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 840).toISOString(),
  },
  {
    id: "04adad90-b0dd-4802-1724-2634d54f940a",
    hash_key: "1c29fd3b8702af0d8d6e23f188f5e8bd16df3afdf38e82df3d8488aef74bef44",
    title: "AI / ML Research Intern - NLP & LLM",
    company: "Sarvam AI",
    role: "Data / AI",
    city: "Bangalore",
    source_url: "https://www.sarvam.ai/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 960).toISOString(),
  },
  {
    id: "f3bebe80-afcc-47f1-0613-1523c43e830b",
    hash_key: "2d3a0e4c9813b01e9e7f34029906f9ce27e04b0ee49f93e04e9599bfa85cfa55",
    title: "Junior QA Automation Engineer",
    company: "BrowserStack",
    role: "QA / Testing",
    city: "Mumbai",
    source_url: "https://www.browserstack.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 1100).toISOString(),
  },
  {
    id: "e2cfcf70-9ebb-46e0-f502-0412b32d720c",
    hash_key: "3e4b1f5da924c12faf804513aa170adf38f15c1ff5a0a4f15fa6aac0b96dfb66",
    title: "Associate Software Engineer - Golang & Microservices",
    company: "Urban Company",
    role: "Backend Engineer",
    city: "Delhi NCR",
    source_url: "https://www.urbancompany.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 1250).toISOString(),
  },
  {
    id: "d1d0d060-8daa-45df-e4f1-f301a21c610d",
    hash_key: "4f5c206eb035d230bf915624bb281bee49026d2006b1b50260b7bbd1ca7e0c77",
    title: "Junior Product Engineer (Remote)",
    company: "Hasura",
    role: "Full Stack",
    city: "Remote",
    source_url: "https://hasura.io/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 1400).toISOString(),
  },
  {
    id: "c0e1e150-7c99-44ce-d3e0-e2f0910b500e",
    hash_key: "506d317fc146e341c0a26735cc392cff5a137e3117c2c61371c8cce2db8f1d88",
    title: "Site Reliability Engineer (Fresher)",
    company: "Zomato",
    role: "DevOps",
    city: "Delhi NCR",
    source_url: "https://www.zomato.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 1600).toISOString(),
  },
  {
    id: "bfa2f240-6b88-43bd-c2df-d1e080fa400f",
    hash_key: "617e4280d257f452d1b37846dd4a3d006b248f4228d3d72482d9ddf3ec902e99",
    title: "Frontend Engineering Intern",
    company: "CoinDCX",
    role: "Frontend Developer",
    city: "Mumbai",
    source_url: "https://coindcx.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 1800).toISOString(),
  },
  {
    id: "aeb30330-5a77-42ac-b1ce-c0df7fea3010",
    hash_key: "728f5391e3680563e2c48957ee5b4e117c35905339e4e83593eaee04fd013faa",
    title: "Software Engineer (Backend - Distributed Systems)",
    company: "Jio Platforms",
    role: "Backend Engineer",
    city: "Hyderabad",
    source_url: "https://careers.jio.com",
    created_at: new Date(Date.now() - 1000 * 60 * 2000).toISOString(),
  },
  {
    id: "9dc41420-4966-419b-a0bd-b0ce6ed92011",
    hash_key: "839064a2f4791674f3d59a68ff6c5f228d46a1644af5f946a4fbf0150e1240bb",
    title: "Machine Learning Intern",
    company: "Fractal Analytics",
    role: "Data / AI",
    city: "Mumbai",
    source_url: "https://fractal.ai/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 2200).toISOString(),
  },
  {
    id: "8cb52510-3855-408a-9fac-a0bd5dc81012",
    hash_key: "94a175b3058a278504e6ab79007d60339e57b2755ba60a57b50c01261f2351cc",
    title: "Junior Backend Developer - Python / FastAPI",
    company: "Klub",
    role: "Backend Engineer",
    city: "Bangalore",
    source_url: "https://klubworks.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 2500).toISOString(),
  },
  {
    id: "7ba63600-2744-4f79-8e9b-9fac4cb70013",
    hash_key: "a5b286c4169b389615f7bc8a118e7144af68c3866cb71b68c61d1237203462dd",
    title: "Full Stack Engineer - Early Career",
    company: "Darwinbox",
    role: "Full Stack",
    city: "Hyderabad",
    source_url: "https://darwinbox.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 2800).toISOString(),
  },
  {
    id: "6a9747f0-1633-4e68-7d8a-8eb03ba6f014",
    hash_key: "b6c397d527ac49a72608cd9b229f8255b079d4977dc82c79d72e2348314573ee",
    title: "Associate React Native Developer",
    company: "Jupiter Money",
    role: "Mobile",
    city: "Bangalore",
    source_url: "https://jupiter.money/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 3100).toISOString(),
  },
  {
    id: "598858e0-0522-4d57-6c79-7da12a95e015",
    hash_key: "c7d4a8e638bd5ab83719deac33a09366c18ae5a88ed93d8ae83f3459425684ff",
    title: "Software Engineer - Quality Assurance",
    company: "TCS",
    role: "QA / Testing",
    city: "Chennai",
    source_url: "https://tcs.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 3400).toISOString(),
  },
  {
    id: "487969d0-f411-4c46-5b68-6c921984d016",
    hash_key: "d8e5b9f749ce6bc9482aefbd44b1a477d29bf6b99fe04e9bf940456a53679500",
    title: "Junior Cloud Engineer (AWS / Azure)",
    company: "Wipro Digital",
    role: "DevOps",
    city: "Pune",
    source_url: "https://wipro.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 3700).toISOString(),
  },
  {
    id: "376a7ac0-e300-4b35-4a57-5b830873c017",
    hash_key: "e9f6ca0850df7c7a593bf0ce55c2b588e3ac07cac05b59ac0a51557bd418be11",
    title: "Data Engineer Trainee - Spark & SQL",
    company: "Mu Sigma",
    role: "Data / AI",
    city: "Bangalore",
    source_url: "https://mu-sigma.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 4000).toISOString(),
  },
  {
    id: "265b8bb0-d2ff-4a24-3946-4a74f762b018",
    hash_key: "fa07db1961e08d8b604c01df66d3c699f4bd18dbd16c60bd1b62668ce529cf22",
    title: "Frontend Developer - Vue.js / Nuxt",
    company: "CleverTap",
    role: "Frontend Developer",
    city: "Mumbai",
    source_url: "https://clevertap.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 4300).toISOString(),
  },
  {
    id: "154c9ca0-c1ee-4913-2835-3965e651a019",
    hash_key: "0b18ec2a72f19e9c715d12e077e4d7aa05ce29ece27d71ce2c73779df63ade33",
    title: "Graduate Software Engineer (C++ / Go)",
    company: "Tower Research Capital",
    role: "Backend Engineer",
    city: "Delhi NCR",
    source_url: "https://tower-research.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 4600).toISOString(),
  },
  {
    id: "043da590-b0dd-4802-1724-2854d5409020",
    hash_key: "1c29fd3b8302af0d826e23f188f5e8bb16df3afdf38e82df3d8488aef74bef44",
    title: "DevSecOps Intern",
    company: "QuickHeal",
    role: "DevOps",
    city: "Pune",
    source_url: "https://quickheal.com/careers",
    created_at: new Date(Date.now() - 1000 * 60 * 5000).toISOString(),
  },
];

// Calculate Levenshtein distance between two strings for typo tolerance.
export function levenshteinDistance(s1: string, s2: string): number {
  const m = s1.length;
  const n = s2.length;
  const dp: number[][] = Array.from({ length: m + 1 }, () => Array(n + 1).fill(0));

  for (let i = 0; i <= m; i++) dp[i][0] = i;
  for (let j = 0; j <= n; j++) dp[0][j] = j;

  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      if (s1[i - 1] === s2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1];
      } else {
        dp[i][j] = 1 + Math.min(dp[i - 1][j], dp[i][j - 1], dp[i - 1][j - 1]);
      }
    }
  }
  return dp[m][n];
}

// Common tech keywords, role misspellings, and Indian city abbreviations
const COMMON_TYPOS: Record<string, string[]> = {
  frotend: ["frontend", "react", "ui"],
  frntend: ["frontend"],
  "front-end": ["frontend"],
  backnd: ["backend", "golang", "node", "java", "python"],
  bakend: ["backend"],
  "back-end": ["backend"],
  fullstak: ["fullstack", "full stack"],
  fulstack: ["fullstack", "full stack"],
  pythn: ["python"],
  pyton: ["python"],
  softwar: ["software", "sde"],
  softwere: ["software"],
  enginer: ["engineer"],
  engneer: ["engineer"],
  bngalore: ["bangalore", "bengaluru"],
  banglore: ["bangalore", "bengaluru"],
  bengaluru: ["bangalore"],
  bnglr: ["bangalore"],
  hydrabad: ["hyderabad"],
  hyderbad: ["hyderabad"],
  mumbay: ["mumbai"],
  bombay: ["mumbai"],
  delhi: ["delhi", "noida", "gurgaon"],
  gurugram: ["gurgaon", "delhi"],
  chenai: ["chennai"],
  remot: ["remote"],
  remotly: ["remote"],
  wfh: ["remote"],
  freshr: ["fresher", "intern", "junior"],
  intren: ["intern", "internship"],
  intrn: ["intern"],
  devloper: ["developer"],
  devlopr: ["developer"],
};

// matchesFuzzy returns true if target text matches the query via substring, dictionary expansion, or Levenshtein distance <= 2.
export function matchesFuzzy(text: string | null | undefined, query: string): boolean {
  if (!text || !query) return false;
  const targetLower = text.toLowerCase().trim();
  const queryLower = query.toLowerCase().trim();

  // 1. Direct substring match
  if (targetLower.includes(queryLower)) return true;

  // 2. Check typo dictionary for full query
  const expansions = COMMON_TYPOS[queryLower] || [];
  for (const exp of expansions) {
    if (targetLower.includes(exp)) return true;
  }

  // 3. Token-level fuzzy comparison
  const targetWords = targetLower.split(/[\s,.\-_/\\()]+/);
  const queryWords = queryLower.split(/\s+/);

  for (const qWord of queryWords) {
    if (qWord.length <= 2) continue; // skip single letters / very short words

    // Check token dictionary expansion
    const tokenExpansions = COMMON_TYPOS[qWord] || [];
    for (const exp of tokenExpansions) {
      if (targetLower.includes(exp)) return true;
    }

    // Levenshtein threshold: 1 edit for 3-4 chars, 2 edits for 5+ chars
    const maxDistance = qWord.length >= 5 ? 2 : 1;

    for (const tWord of targetWords) {
      if (tWord.length <= 2) continue;
      if (Math.abs(tWord.length - qWord.length) > maxDistance) continue;
      if (levenshteinDistance(tWord, qWord) <= maxDistance) {
        return true;
      }
    }
  }

  return false;
}

export async function fetchJobs(
  params: JobFilterParams = {}
): Promise<{ response: PaginatedJobsResponse; isFallback: boolean }> {
  const { page = 1, limit = 20, role, city, search } = params;

  const query = new URLSearchParams();
  query.set("page", page.toString());
  query.set("limit", limit.toString());

  if (role && role !== "All Roles") {
    query.set("role", role);
  }
  if (city && city !== "All Cities") {
    query.set("city", city);
  }
  if (search && search.trim() !== "") {
    query.set("search", search.trim());
  }

  const endpoint = `${API_BASE_URL}/api/jobs?${query.toString()}`;

  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 3500);

    const res = await fetch(endpoint, {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    if (res.ok) {
      const data: PaginatedJobsResponse = await res.json();
      let filteredData = data.data || [];
      if (search && search.trim() !== "") {
        filteredData = filteredData.filter(
          (j) =>
            matchesFuzzy(j.title, search) ||
            matchesFuzzy(j.company, search) ||
            matchesFuzzy(j.role, search) ||
            matchesFuzzy(j.city, search)
        );
      }

      return {
        response: {
          data: filteredData,
          pagination: data.pagination,
        },
        isFallback: false,
      };
    }
  } catch {
    // Backend unreachable or offline, graceful fallback
  }

  // Fallback to rich curated local dataset
  let results = [...SAMPLE_JOBS];

  if (role && role !== "All Roles") {
    results = results.filter((j) =>
      j.role?.toLowerCase().includes(role.toLowerCase())
    );
  }

  if (city && city !== "All Cities") {
    results = results.filter((j) =>
      j.city?.toLowerCase().includes(city.toLowerCase())
    );
  }

  if (search && search.trim() !== "") {
    results = results.filter(
      (j) =>
        matchesFuzzy(j.title, search) ||
        matchesFuzzy(j.company, search) ||
        matchesFuzzy(j.role, search) ||
        matchesFuzzy(j.city, search)
    );
  }

  const total = results.length;
  const totalPages = Math.ceil(total / limit) || 1;
  const safePage = Math.min(Math.max(1, page), totalPages);
  const start = (safePage - 1) * limit;
  const end = start + limit;
  const paginatedData = results.slice(start, end);

  return {
    response: {
      data: paginatedData,
      pagination: {
        total,
        page: safePage,
        limit,
        total_pages: totalPages,
        has_next: safePage < totalPages,
        has_prev: safePage > 1,
      },
    },
    isFallback: true,
  };
}
