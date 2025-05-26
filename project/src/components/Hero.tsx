import React from 'react';
import { ArrowRight } from 'lucide-react';

const Hero = () => {
  return (
    <section 
      id="home" 
      className="relative pt-28 pb-24 min-h-screen flex items-center bg-gradient-to-br from-slate-900 to-indigo-900"
    >
      <div className="absolute inset-0 bg-[url('https://images.pexels.com/photos/1287142/pexels-photo-1287142.jpeg?auto=compress&cs=tinysrgb&w=1920')] bg-cover bg-center opacity-20"></div>
      <div className="absolute inset-0 bg-gradient-to-b from-slate-900/90 to-indigo-900/80"></div>
      
      <div className="relative container mx-auto px-4 sm:px-6 lg:px-8 z-10">
        <div className="max-w-4xl mx-auto text-center">
          <div className="inline-block px-3 py-1 mb-6 rounded-full bg-indigo-600/20 border border-indigo-500/30">
            <p className="text-xs sm:text-sm font-medium text-indigo-300">
              Domain testing with dAppling & 3DNS
            </p>
          </div>
          
          <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold text-white mb-6 leading-tight tracking-tight">
            Your Decentralized <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-sky-400">Web3 Domain</span> is Ready
          </h1>
          
          <p className="text-lg sm:text-xl text-slate-300 mb-8 max-w-2xl mx-auto">
            Congratulations! Your domain is now connected to a decentralized hosting service. This is a test page to verify your domain is working correctly.
          </p>
          
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <button className="w-full sm:w-auto bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-3 rounded-lg text-base font-medium transition-all transform hover:scale-105 hover:shadow-lg flex items-center justify-center">
              <span>Configure Domain</span>
              <ArrowRight className="ml-2 h-5 w-5" />
            </button>
            <button className="w-full sm:w-auto bg-white/10 hover:bg-white/20 text-white border border-white/20 px-8 py-3 rounded-lg text-base font-medium transition-all hover:shadow-lg">
              Learn More
            </button>
          </div>
        </div>
      </div>

      <div className="absolute bottom-6 left-0 right-0 flex justify-center animate-bounce">
        <a href="#features" className="text-white/60 hover:text-white">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </a>
      </div>
    </section>
  );
};

export default Hero;