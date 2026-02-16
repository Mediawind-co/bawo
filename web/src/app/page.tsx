import Image from "next/image";
import Link from "next/link";
import { FireIcon, BookIcon, TrophyIcon, UsersIcon, CheckIcon, PlayIcon } from "@/components/icons";

const features = [
  {
    icon: BookIcon,
    title: "Interactive Lessons",
    description: "Learn through listening, speaking, and interactive exercises designed for real-world conversations.",
  },
  {
    icon: FireIcon,
    title: "Daily Streaks",
    description: "Build consistent habits with daily goals and streak tracking to keep you motivated.",
  },
  {
    icon: TrophyIcon,
    title: "Track Progress",
    description: "See your improvement over time with detailed statistics and achievement milestones.",
  },
  {
    icon: UsersIcon,
    title: "Community",
    description: "Join thousands of learners on their journey to master African languages.",
  },
];

const languages = [
  { name: "Yoruba", emoji: "🇳🇬", speakers: "45M+ speakers" },
  { name: "Igbo", emoji: "🇳🇬", speakers: "30M+ speakers" },
  { name: "Hausa", emoji: "🇳🇬", speakers: "70M+ speakers" },
];

const benefits = [
  "Bite-sized lessons that fit your schedule",
  "Audio from native speakers",
  "Progress tracking and XP rewards",
  "Review system for mistakes",
  "Works on any device",
];

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white">
      {/* Navigation */}
      <nav className="fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link href="/" className="flex items-center gap-2">
              <Image src="/logo.png" alt="Bawo" width={120} height={32} className="h-8 w-auto" />
            </Link>
            <div className="hidden md:flex items-center gap-8">
              <a href="#features" className="text-teal hover:text-mint transition-colors">Features</a>
              <a href="#languages" className="text-teal hover:text-mint transition-colors">Languages</a>
              <a href="#about" className="text-teal hover:text-mint transition-colors">About</a>
            </div>
            <div className="flex items-center gap-4">
              <Link
                href="/login"
                className="text-teal hover:text-mint transition-colors font-medium"
              >
                Log in
              </Link>
              <Link
                href="/login"
                className="bg-mint hover:bg-primary-dark text-teal font-semibold px-5 py-2.5 rounded-full transition-colors"
              >
                Get Started
              </Link>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <div className="text-center max-w-3xl mx-auto">
            <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold text-teal leading-tight">
              Learn African Languages
              <span className="text-gradient block mt-2">The Fun Way</span>
            </h1>
            <p className="mt-6 text-xl text-gray-600 max-w-2xl mx-auto">
              Master Yoruba, Igbo, Hausa and more with interactive lessons, native audio, and personalized learning paths.
            </p>
            <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link
                href="/login"
                className="w-full sm:w-auto bg-mint hover:bg-primary-dark text-teal font-semibold px-8 py-4 rounded-full text-lg transition-colors flex items-center justify-center gap-2"
              >
                <PlayIcon className="w-5 h-5" />
                Start Learning Free
              </Link>
              <a
                href="#features"
                className="w-full sm:w-auto border-2 border-teal text-teal hover:bg-teal hover:text-white font-semibold px-8 py-4 rounded-full text-lg transition-colors"
              >
                See How It Works
              </a>
            </div>
          </div>

          {/* Hero illustration/stats */}
          <div className="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-4xl mx-auto">
            {languages.map((lang) => (
              <div
                key={lang.name}
                className="bg-gradient-to-br from-mint/10 to-purple/5 rounded-2xl p-6 text-center border border-mint/20"
              >
                <span className="text-5xl">{lang.emoji}</span>
                <h3 className="mt-4 text-xl font-bold text-teal">{lang.name}</h3>
                <p className="mt-1 text-gray-600">{lang.speakers}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h2 className="text-3xl sm:text-4xl font-bold text-teal">Why Learn with Bawo?</h2>
            <p className="mt-4 text-xl text-gray-600 max-w-2xl mx-auto">
              Our platform is designed to make learning African languages accessible, engaging, and effective.
            </p>
          </div>
          <div className="mt-16 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="bg-white rounded-2xl p-6 shadow-sm hover:shadow-md transition-shadow border border-gray-100"
              >
                <div className="w-12 h-12 bg-mint/20 rounded-xl flex items-center justify-center">
                  <feature.icon className="w-6 h-6 text-teal" />
                </div>
                <h3 className="mt-4 text-xl font-semibold text-teal">{feature.title}</h3>
                <p className="mt-2 text-gray-600">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Languages Section */}
      <section id="languages" className="py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            <div>
              <h2 className="text-3xl sm:text-4xl font-bold text-teal">
                Languages That Connect You to Culture
              </h2>
              <p className="mt-4 text-xl text-gray-600">
                Learn languages spoken by over 145 million people across Africa. Connect with your heritage or explore new cultures.
              </p>
              <ul className="mt-8 space-y-4">
                {benefits.map((benefit) => (
                  <li key={benefit} className="flex items-center gap-3">
                    <div className="w-6 h-6 bg-mint rounded-full flex items-center justify-center flex-shrink-0">
                      <CheckIcon className="w-4 h-4 text-teal" />
                    </div>
                    <span className="text-gray-700">{benefit}</span>
                  </li>
                ))}
              </ul>
              <Link
                href="/login"
                className="mt-8 inline-block bg-purple hover:bg-accent text-white font-semibold px-8 py-4 rounded-full text-lg transition-colors"
              >
                Start Your Journey
              </Link>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-4">
                <div className="bg-gradient-to-br from-mint to-mint/70 rounded-2xl p-6 text-teal">
                  <p className="text-4xl font-bold">500+</p>
                  <p className="mt-1 font-medium">Lessons</p>
                </div>
                <div className="bg-gradient-to-br from-purple to-indigo rounded-2xl p-6 text-white">
                  <p className="text-4xl font-bold">1000+</p>
                  <p className="mt-1 font-medium">Audio Clips</p>
                </div>
              </div>
              <div className="space-y-4 mt-8">
                <div className="bg-gradient-to-br from-teal to-teal/80 rounded-2xl p-6 text-white">
                  <p className="text-4xl font-bold">50k+</p>
                  <p className="mt-1 font-medium">Active Learners</p>
                </div>
                <div className="bg-gradient-to-br from-indigo to-purple rounded-2xl p-6 text-white">
                  <p className="text-4xl font-bold">95%</p>
                  <p className="mt-1 font-medium">Satisfaction</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 gradient-primary">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl sm:text-4xl font-bold text-teal">
            Ready to Start Learning?
          </h2>
          <p className="mt-4 text-xl text-teal/80">
            Join thousands of learners mastering African languages today.
          </p>
          <Link
            href="/login"
            className="mt-8 inline-block bg-teal hover:bg-teal/90 text-white font-semibold px-10 py-4 rounded-full text-lg transition-colors"
          >
            Get Started for Free
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer id="about" className="bg-teal text-white py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-12">
            <div className="col-span-1 md:col-span-2">
              <Image src="/logo.png" alt="Bawo" width={120} height={32} className="h-8 w-auto brightness-0 invert" />
              <p className="mt-4 text-white/80 max-w-md">
                Bawo is dedicated to preserving and spreading African languages through modern, accessible language learning technology.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-lg">Learn</h4>
              <ul className="mt-4 space-y-2 text-white/80">
                <li><a href="#" className="hover:text-mint transition-colors">Yoruba</a></li>
                <li><a href="#" className="hover:text-mint transition-colors">Igbo</a></li>
                <li><a href="#" className="hover:text-mint transition-colors">Hausa</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-lg">Company</h4>
              <ul className="mt-4 space-y-2 text-white/80">
                <li><a href="#" className="hover:text-mint transition-colors">About Us</a></li>
                <li><a href="#" className="hover:text-mint transition-colors">Privacy</a></li>
                <li><a href="#" className="hover:text-mint transition-colors">Terms</a></li>
              </ul>
            </div>
          </div>
          <div className="mt-12 pt-8 border-t border-white/20 text-center text-white/60">
            <p>&copy; {new Date().getFullYear()} Bawo. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
