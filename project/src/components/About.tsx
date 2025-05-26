import React from 'react';

const About = () => {
  return (
    <section id="about" className="py-24 bg-slate-50">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col lg:flex-row items-center gap-12">
          <div className="w-full lg:w-1/2">
            <div className="relative">
              <div className="absolute -top-4 -left-4 w-24 h-24 bg-indigo-500/10 rounded-lg z-0"></div>
              <div className="absolute -bottom-4 -right-4 w-24 h-24 bg-indigo-500/10 rounded-lg z-0"></div>
              <div className="relative overflow-hidden rounded-xl shadow-xl z-10">
                <img 
                  src="https://images.pexels.com/photos/8092507/pexels-photo-8092507.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2" 
                  alt="Web3 Domain Visualization" 
                  className="w-full h-auto"
                />
              </div>
            </div>
          </div>
          
          <div className="w-full lg:w-1/2">
            <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-6">
              What is dAppling & 3DNS?
            </h2>
            
            <div className="space-y-6 text-gray-600">
              <p>
                <strong>dAppling</strong> is a decentralized hosting platform that allows you to deploy and serve web applications without relying on traditional centralized hosting providers.
              </p>
              
              <p>
                <strong>3DNS (Decentralized Domain Name System)</strong> connects traditional domain names to decentralized content, bridging the gap between Web2 and Web3 technologies.
              </p>
              
              <p>
                By connecting your domain to these services, you're taking a step into the future of the internet - one that's more secure, resistant to censorship, and puts control back in the hands of users.
              </p>
              
              <div className="pt-4">
                <button className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-lg text-base font-medium transition-all transform hover:scale-105 hover:shadow-lg">
                  Learn More About Web3 Domains
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default About;